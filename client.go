package cfg

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/10antz-inc/pf-go-cfg/decoder"
	"github.com/10antz-inc/pf-go-cfg/encoder"
	"github.com/10antz-inc/pf-go-cfg/option"
	"github.com/10antz-inc/pf-go-cfg/pubsub"
	"github.com/10antz-inc/pf-go-cfg/store"
	store_option "github.com/10antz-inc/pf-go-cfg/store/option"
	"github.com/10antz-inc/pf-go-ers"
)

type Client interface {
	Get(ctx context.Context) (any, error)
	Set(ctx context.Context, value any) error
	Encode(value any) ([]byte, error)
	Decode(value []byte) (any, error)
	Available() bool
	Close() error
}

type client struct {
	mu         map[store.Store]*sync.RWMutex
	msg        *message
	origin     store.Store
	pubsub     pubsub.PubSub
	cacheStore store.Store
	available  bool
	encoder    encoder.Encoder
	decoder    decoder.Decoder
	option     option.Client
}

var _ Client = (*client)(nil)

func NewClient(ctx context.Context, msg any, origin store.Store, pubsub pubsub.PubSub, options ...option.ClientOption) (Client, error) {
	c := &client{
		msg:       newMessage(msg),
		origin:    origin,
		pubsub:    pubsub,
		available: true,
	}

	for _, option := range options {
		option(&c.option)
	}

	if v := c.option.CacheStore; v != nil {
		c.cacheStore = v
	} else {
		c.cacheStore = store.NewMemory(
			store_option.WithDefaultExpiration(time.Duration(5*time.Minute)),
			store_option.WithCleanupInterval(60*time.Minute),
		)
	}

	c.mu = map[store.Store]*sync.RWMutex{}
	c.mu[c.origin] = &sync.RWMutex{}
	c.mu[c.cacheStore] = &sync.RWMutex{}

	if v := c.option.Encoder; v != nil {
		c.encoder = v
	} else {
		c.encoder = encoder.NewJSON()
	}

	if v := c.option.Decoder; v != nil {
		c.decoder = v
	} else {
		c.decoder = decoder.NewJSON()
	}

	go func() {
		ctx := context.Background()
		err := pubsub.Subscribe(ctx, func(ctx context.Context, _ []byte) error {
			// 更新が発生したためキャッシュをクリアする
			return ers.W(c.del(ctx, c.cacheStore))
		})
		if err != nil {
			c.available = false
		}
	}()

	return c, nil
}

func (c *client) Get(ctx context.Context) (val any, err error) {
	val, err = c.get(ctx, c.cacheStore)
	if err != nil {
		return nil, ers.W(err)
	}
	if val != nil {
		return val, nil
	}

	val, err = c.get(ctx, c.origin)
	if err != nil {
		return nil, ers.W(err)
	}
	if val != nil {
		if err := c.set(ctx, c.cacheStore, val); err != nil {
			return nil, ers.W(err)
		}
		return val, nil
	}

	return nil, nil
}

func (c *client) Set(ctx context.Context, value any) error {
	if err := c.set(ctx, c.cacheStore, value); err != nil {
		return ers.W(err)
	}

	if err := c.set(ctx, c.origin, value); err != nil {
		return ers.W(err)
	}

	if err := c.pubsub.Publish(ctx, []byte("update")); err != nil {
		return ers.W(err)
	}

	return nil
}

func (c *client) Encode(val any) ([]byte, error) {
	bytes, err := c.encoder.Encode(val)
	if err != nil {
		return nil, ers.W(err)
	}
	return bytes, nil
}

func (c *client) Decode(bytes []byte) (any, error) {
	val := c.msg.new()
	if err := c.decoder.Decode(bytes, val); err != nil {
		return nil, ers.W(err)
	}

	// new で生成されるものはポインタなので参照先を返却する
	return reflect.ValueOf(val).Elem().Interface(), nil
}

func (c *client) Available() bool {
	return c.available
}

func (c *client) Close() error {
	if err := c.pubsub.Close(); err != nil {
		return ers.W(err)
	}
	return nil
}

func (c *client) get(ctx context.Context, store store.Store) (any, error) {
	c.mu[store].RLock()
	defer c.mu[store].RUnlock()

	bytes, err := store.Get(ctx, c.msg.Name())
	if err != nil {
		return nil, ers.W(err)
	} else if bytes == nil {
		return nil, nil
	}

	val, err := c.Decode(bytes)
	if err != nil {
		return nil, ers.W(err)
	} else if c.msg.Name() != newMessage(val).Name() {
		return nil, ers.ErrInvalidArgument.New("message type is mismatch")
	}

	return val, nil
}

func (c *client) set(ctx context.Context, store store.Store, value any) error {
	c.mu[store].Lock()
	defer c.mu[store].Unlock()

	if c.msg.Name() != newMessage(value).Name() {
		return ers.ErrInvalidArgument.New("message type is mismatch")
	}

	options := []store_option.CacheOption{}
	if v := c.option.Expiration; v != 0 {
		options = append(options, store_option.WithExpiration(v))
	}

	bytes, err := c.Encode(value)
	if err != nil {
		return ers.W(err)
	}

	if err := store.Set(ctx, c.msg.Name(), bytes, options...); err != nil {
		return ers.W(err)
	}

	return nil
}

func (c *client) del(ctx context.Context, store store.Store) error {
	c.mu[store].Lock()
	defer c.mu[store].Unlock()

	if err := store.Del(ctx, c.msg.Name()); err != nil {
		return ers.W(err)
	}

	return nil
}

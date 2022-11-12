package cfg

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/tys-muta/go-cfg/decoder"
	"github.com/tys-muta/go-cfg/encoder"
	"github.com/tys-muta/go-cfg/errors"
	"github.com/tys-muta/go-cfg/option"
	"github.com/tys-muta/go-cfg/pubsub"
	"github.com/tys-muta/go-cfg/store"
	s_option "github.com/tys-muta/go-cfg/store/option"
	"github.com/tys-muta/go-ers"
	"github.com/tys-muta/go-opt"
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
	mu        map[store.Store]*sync.RWMutex
	msg       *message
	origin    store.Store
	cache     store.Store
	pubsub    pubsub.PubSub
	encoder   encoder.Encoder
	decoder   decoder.Decoder
	options   option.ClientOptions
	available bool
}

var _ Client = (*client)(nil)

func NewClient(ctx context.Context, msg any, origin store.Store, pubsub pubsub.PubSub, options ...opt.Option) (Client, error) {
	c := &client{
		msg:       newMessage(msg),
		origin:    origin,
		pubsub:    pubsub,
		available: true,
	}

	if err := opt.Reflect(&c.options, options...); err != nil {
		return nil, ers.W(err)
	}

	if v := c.options.Cache; v != nil {
		c.cache = v
	} else {
		c.cache = store.NewMemory(
			s_option.WithDefaultExpiration(5*time.Minute),
			s_option.WithCleanupInterval(60*time.Minute),
		)
	}

	c.mu = map[store.Store]*sync.RWMutex{}
	c.mu[c.origin] = &sync.RWMutex{}
	c.mu[c.cache] = &sync.RWMutex{}

	if v := c.options.Encoder; v != nil {
		c.encoder = v
	} else {
		c.encoder = encoder.NewJSON()
	}

	if v := c.options.Decoder; v != nil {
		c.decoder = v
	} else {
		c.decoder = decoder.NewJSON()
	}

	go func() {
		if err := pubsub.Subscribe(ctx, func(ctx context.Context, _ []byte) error {
			if err := c.del(ctx, c.cache); err != nil {
				return ers.W(err)
			}
			return nil
		}); err != nil {
			if ers.Is(err, errors.ErrNotFoundSubscription) {
				c.available = false
			}
		}
	}()

	return c, nil
}

func (c *client) Get(ctx context.Context) (any, error) {
	if v, err := c.get(ctx, c.cache); err != nil {
		return nil, ers.W(err)
	} else if v != nil {
		return v, nil
	}

	if v, err := c.get(ctx, c.origin); err != nil {
		return nil, ers.W(err)
	} else if v != nil {
		if err := c.set(ctx, c.cache, v); err != nil {
			return nil, ers.W(err)
		}
		return v, nil
	}

	return nil, nil
}

func (c *client) Set(ctx context.Context, value any) error {
	if err := c.set(ctx, c.cache, value); err != nil {
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

func (c *client) Encode(value any) ([]byte, error) {
	if v, err := c.encoder.Encode(value); err != nil {
		return nil, ers.W(err)
	} else {
		return v, nil
	}
}

func (c *client) Decode(value []byte) (any, error) {
	v := c.msg.new()
	if err := c.decoder.Decode(value, v); err != nil {
		return nil, ers.W(err)
	} else {
		// new で生成されるものはポインタなので参照先を返却する
		v = reflect.ValueOf(v).Elem().Interface()
		return v, nil
	}
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

	if v, err := store.Get(ctx, c.msg.Name()); err != nil {
		return nil, ers.W(err)
	} else if v == nil {
		return nil, nil
	} else if v, err := c.Decode(v); err != nil {
		return nil, ers.W(err)
	} else if c.msg.Name() != newMessage(v).Name() {
		return nil, ers.ErrInvalidArgument.New("message type is mismatch")
	} else {
		return v, nil
	}
}

func (c *client) set(ctx context.Context, store store.Store, value any) error {
	c.mu[store].Lock()
	defer c.mu[store].Unlock()

	if c.msg.Name() != newMessage(value).Name() {
		return ers.ErrInvalidArgument.New("message type is mismatch")
	}

	options := []opt.Option{}
	if v := c.options.Expiration; v != nil {
		options = append(options, s_option.WithExpiration(time.Duration(*v)))
	}

	if v, err := c.Encode(value); err != nil {
		return ers.W(err)
	} else if err := store.Set(ctx, c.msg.Name(), v, options...); err != nil {
		return ers.W(err)
	} else {
		return nil
	}
}

func (c *client) del(ctx context.Context, store store.Store) error {
	c.mu[store].Lock()
	defer c.mu[store].Unlock()

	if err := store.Del(ctx, c.msg.Name()); err != nil {
		return ers.W(err)
	} else {
		return nil
	}
}

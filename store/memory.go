package store

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
	s_option "github.com/tys-muta/go-cfg/store/option"
	"github.com/tys-muta/go-ers"
	"github.com/tys-muta/go-opt"
)

type memory struct {
	client *cache.Cache
}

var _ Store = (*memory)(nil)

func NewMemory(options ...opt.Option) (Store, error) {
	s := &memory{}

	o := &s_option.MemoryOptions{}
	if err := opt.Reflect(o, options...); err != nil {
		return nil, ers.W(err)
	}

	var defaultExpiration time.Duration
	if v := o.DefaultExpiration; v != nil {
		defaultExpiration = time.Duration(*v)
	}
	var cleanupInterval time.Duration
	if v := o.CleanupInterval; v != nil {
		cleanupInterval = time.Duration(*v)
	}

	s.client = cache.New(defaultExpiration, cleanupInterval)

	return s, nil
}

func (s *memory) Get(ctx context.Context, key string) ([]byte, error) {
	if v, ok := s.client.Get(key); !ok {
		return nil, nil
	} else if v, ok := v.([]byte); !ok {
		return nil, ers.ErrInternal.New("failed to type assertion")
	} else {
		return v, nil
	}
}

func (s *memory) Set(ctx context.Context, key string, value []byte, options ...opt.Option) error {
	o := &s_option.CacheOptions{}
	if err := opt.Reflect(o, options...); err != nil {
		return ers.W(err)
	}

	var duration time.Duration = cache.DefaultExpiration
	if v := o.Expiration; v != nil {
		duration = time.Duration(*v)
	}

	s.client.Set(key, value, duration)

	return nil
}

func (s *memory) Del(ctx context.Context, key string) error {
	s.client.Delete(key)

	return nil
}

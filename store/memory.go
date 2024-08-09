package store

import (
	"context"

	"github.com/10antz-inc/pf-go-cfg/store/option"
	"github.com/10antz-inc/pf-go-ers"
	"github.com/patrickmn/go-cache"
)

type memory struct {
	client *cache.Cache
}

var _ Store = (*memory)(nil)

func NewMemory(options ...option.MemoryOption) *memory {
	s := &memory{}

	o := option.MemoryOptions{}
	for _, option := range options {
		option(&o)
	}

	s.client = cache.New(o.DefaultExpiration, o.CleanupInterval)

	return s
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

func (s *memory) Set(ctx context.Context, key string, value []byte, options ...option.CacheOption) error {
	o := option.CacheOptions{}
	for _, option := range options {
		option(&o)
	}

	s.client.Set(key, value, o.Expiration)

	return nil
}

func (s *memory) Del(ctx context.Context, key string) error {
	s.client.Delete(key)

	return nil
}

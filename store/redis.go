package store

import (
	"context"

	"github.com/10antz-inc/pf-go-cfg/store/option"
	p_redis "github.com/go-redis/redis"
)

// 必要になった時に実装
type redis struct {
	client *p_redis.Client
}

var _ Store = (*redis)(nil)

func NewRedis(client *p_redis.Client) (Store, error) {
	s := &redis{client: client}
	return s, nil
}

func (s *redis) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}

func (s *redis) Set(ctx context.Context, key string, value []byte, options ...option.CacheOption) error {
	return nil
}

func (s *redis) Del(ctx context.Context, key string) error {
	return nil
}

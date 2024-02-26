package store

import (
	"context"

	"github.com/tys-muta/go-cfg/store/option"
)

type Store interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, bytes []byte, options ...option.CacheOption) error
	Del(ctx context.Context, key string) error
}

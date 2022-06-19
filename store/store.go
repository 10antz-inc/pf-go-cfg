package store

import (
	"context"

	"github.com/tys-muta/go-opt"
)

type Store interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, options ...opt.Option) error
	Del(ctx context.Context, key string) error
}

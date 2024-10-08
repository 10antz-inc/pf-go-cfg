package store

import (
	"context"
	p_sql "database/sql"

	"github.com/10antz-inc/pf-go-cfg/store/option"
)

// 必要になった時に実装
type sql struct {
	client *p_sql.DB
}

var _ Store = (*sql)(nil)

func NewSQL(client *p_sql.DB) (Store, error) {
	s := &sql{client: client}
	return s, nil
}

func (s *sql) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}

func (s *sql) Set(ctx context.Context, key string, value []byte, options ...option.CacheOption) error {
	return nil
}

func (s *sql) Del(ctx context.Context, key string) error {
	return nil
}

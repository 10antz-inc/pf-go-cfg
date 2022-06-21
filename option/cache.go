package option

import (
	"github.com/tys-muta/go-cfg/store"
	"github.com/tys-muta/go-ers"
	"github.com/tys-muta/go-opt"
)

type cache struct {
	store.Store
}

var _ opt.Option = (*cache)(nil)

func WithCache(v store.Store) opt.Option {
	return cache{Store: v}
}

func (o cache) Validate() error {
	if o.Store == nil {
		return ers.ErrInvalidArgument.New("store is nil")
	}
	return nil
}

func (o cache) Apply(options any) {
	switch v := options.(type) {
	case *ClientOptions:
		v.Cache = &o
	}
}

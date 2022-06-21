package option

import (
	"time"

	"github.com/tys-muta/go-ers"
	"github.com/tys-muta/go-opt"
)

type expiration time.Duration

var _ opt.Option = (*expiration)(nil)

func WithExpiration(v time.Duration) opt.Option {
	return expiration(v)
}

func (o expiration) Validate() error {
	if o <= 0 {
		return ers.ErrInvalidArgument.New("expiration is invalid")
	}
	return nil
}

func (o expiration) Apply(options any) {
	switch v := options.(type) {
	case *CacheOptions:
		v.Expiration = &o
	}
}

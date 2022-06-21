package option

import (
	"time"

	"github.com/tys-muta/go-ers"
	"github.com/tys-muta/go-opt"
)

type defaultExpiration time.Duration

var _ opt.Option = (*defaultExpiration)(nil)

func WithDefaultExpiration(v time.Duration) opt.Option {
	return defaultExpiration(v)
}

func (o defaultExpiration) Validate() error {
	if o <= 0 {
		return ers.ErrInvalidArgument.New("default expiration is invalid")
	}
	return nil
}

func (o defaultExpiration) Apply(options any) {
	switch v := options.(type) {
	case *MemoryOptions:
		v.DefaultExpiration = &o
	}
}

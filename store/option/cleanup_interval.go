package option

import (
	"time"

	"github.com/tys-muta/go-ers"
	"github.com/tys-muta/go-opt"
)

type cleanupInterval time.Duration

var _ opt.Option = (*cleanupInterval)(nil)

func WithCleanupInterval(v time.Duration) opt.Option {
	return cleanupInterval(v)
}

func (o cleanupInterval) Validate() error {
	if o <= 0 {
		return ers.ErrInvalidArgument.New("cleanup interval is invalid")
	}
	return nil
}

func (o cleanupInterval) Apply(options any) {
	switch v := options.(type) {
	case *MemoryOptions:
		v.CleanupInterval = &o
	}
}

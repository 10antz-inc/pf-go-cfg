package option

import (
	p_encoder "github.com/tys-muta/go-cfg/encoder"
	"github.com/tys-muta/go-ers"
	"github.com/tys-muta/go-opt"
)

type encoder struct {
	p_encoder.Encoder
}

var _ opt.Option = (*encoder)(nil)

func WithEncoder(v p_encoder.Encoder) opt.Option {
	return encoder{Encoder: v}
}

func (o encoder) Validate() error {
	if o.Encoder == nil {
		return ers.ErrInvalidArgument.New("encoder is nil")
	}
	return nil
}

func (o encoder) Apply(options any) {
	switch v := options.(type) {
	case *ClientOptions:
		v.Encoder = &o
	}
}

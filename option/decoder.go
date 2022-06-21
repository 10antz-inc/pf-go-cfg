package option

import (
	p_decoder "github.com/tys-muta/go-cfg/decoder"
	"github.com/tys-muta/go-ers"
	"github.com/tys-muta/go-opt"
)

type decoder struct {
	p_decoder.Decoder
}

var _ opt.Option = (*decoder)(nil)

func WithDecoder(v p_decoder.Decoder) opt.Option {
	return decoder{Decoder: v}
}

func (o decoder) Validate() error {
	if o.Decoder == nil {
		return ers.ErrInvalidArgument.New("decoder is nil")
	}
	return nil
}

func (o decoder) Apply(options any) {
	switch v := options.(type) {
	case *ClientOptions:
		v.Decoder = &o
	}
}

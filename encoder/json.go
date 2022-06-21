package encoder

import (
	p_json "encoding/json"
	"fmt"

	"github.com/tys-muta/go-ers"
)

type json struct{}

var _ Encoder = (*json)(nil)

func NewJSON() Encoder {
	e := &json{}
	return e
}

func (e *json) Encode(v any) ([]byte, error) {
	if v, err := p_json.Marshal(v); err != nil {
		return nil, ers.ErrInternal.New(fmt.Sprintf("failed to marshal: %#v", v))
	} else {
		return v, nil
	}
}

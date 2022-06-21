package decoder

import (
	p_json "encoding/json"
	"fmt"

	"github.com/tys-muta/go-ers"
)

type json struct{}

var _ Decoder = (*json)(nil)

func NewJSON() Decoder {
	d := &json{}
	return d
}

func (d *json) Decode(data []byte, v any) error {
	if err := p_json.Unmarshal(data, v); err != nil {
		return ers.ErrInternal.New(fmt.Sprintf("failed to unmarshal: %t <- %s", v, data))
	}
	return nil
}

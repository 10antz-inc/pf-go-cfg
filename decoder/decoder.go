package decoder

type Decoder interface {
	Decode(bytes []byte, v any) error
}

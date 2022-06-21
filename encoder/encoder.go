package encoder

type Encoder interface {
	Encode(v any) ([]byte, error)
}

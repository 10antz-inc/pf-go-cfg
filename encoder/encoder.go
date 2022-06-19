package encoder

type Encoder interface {
	Encode(v interface{}) ([]byte, error)
}

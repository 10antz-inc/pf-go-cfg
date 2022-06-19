package option

type ClientOptions struct {
	Cache      *cache
	Encoder    *encoder
	Decoder    *decoder
	Expiration *expiration
}

package option

import (
	"time"

	"github.com/tys-muta/go-cfg/decoder"
	"github.com/tys-muta/go-cfg/encoder"
	"github.com/tys-muta/go-cfg/store"
)

type Client struct {
	CacheStore store.Store
	Encoder    encoder.Encoder
	Decoder    decoder.Decoder
	Expiration time.Duration
}

type ClientOption func(options *Client)

func WithCacheStore(v store.Store) ClientOption {
	return func(options *Client) {
		options.CacheStore = v
	}
}

func WithEncoder(v encoder.Encoder) ClientOption {
	return func(options *Client) {
		options.Encoder = v
	}
}

func WithDecoder(v decoder.Decoder) ClientOption {
	return func(options *Client) {
		options.Decoder = v
	}
}

func WithExpiration(v time.Duration) ClientOption {
	return func(options *Client) {
		options.Expiration = v
	}
}

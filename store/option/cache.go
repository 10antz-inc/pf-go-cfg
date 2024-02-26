package option

import "time"

type CacheOptions struct {
	Expiration time.Duration
}

type CacheOption func(options *CacheOptions)

func WithExpiration(v time.Duration) CacheOption {
	return func(options *CacheOptions) {
		options.Expiration = v
	}
}

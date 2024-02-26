package option

import "time"

type MemoryOptions struct {
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
}

type MemoryOption func(options *MemoryOptions)

func WithDefaultExpiration(v time.Duration) MemoryOption {
	return func(options *MemoryOptions) {
		options.DefaultExpiration = v
	}
}

func WithCleanupInterval(v time.Duration) MemoryOption {
	return func(options *MemoryOptions) {
		options.CleanupInterval = v
	}
}

package pubsub

type Closer interface {
	Close() error
}

package pubsub

type PubSub interface {
	Publisher
	Subscriber
	Closer
}

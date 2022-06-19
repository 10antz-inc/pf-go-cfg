package pubsub

import "context"

type Subscriber interface {
	Subscribe(ctx context.Context, subFunc SubscribeFunc) error
}

type SubscribeFunc func(ctx context.Context, msg []byte) error

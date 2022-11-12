package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/tys-muta/go-cfg/errors"
	"github.com/tys-muta/go-ers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type cloudPB struct {
	client         pubsub.Client
	topicID        string
	subscriptionID string
}

var _ PubSub = (*cloudPB)(nil)

func NewCloudPubSub(client pubsub.Client, topicID string, subscriptionID string) *cloudPB {
	return &cloudPB{client: client, topicID: topicID, subscriptionID: subscriptionID}
}

func (p *cloudPB) Publish(ctx context.Context, msg []byte) error {
	if _, err := p.client.Topic(p.topicID).Publish(ctx, &pubsub.Message{Data: msg}).Get(ctx); err != nil {
		return ers.ErrInternal.New(err)
	}

	return nil
}

func (p *cloudPB) Subscribe(ctx context.Context, subFunc SubscribeFunc) error {
	if err := p.client.Subscription(p.subscriptionID).Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		if err := subFunc(ctx, msg.Data); err != nil {
			return
		}
		msg.Ack()
	}); err != nil {
		v, ok := err.(interface{ GRPCStatus() *status.Status })
		if ok && v.GRPCStatus().Code() == codes.NotFound {
			return errors.ErrNotFoundSubscription
		}
		return ers.ErrInternal.New(err)
	}

	return nil
}

func (p *cloudPB) Close() error {
	if err := p.client.Close(); err != nil {
		return ers.W(err)
	}

	return nil
}

package pubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"
	"github.com/tys-muta/go-ers"
)

type cloudPB struct {
	client  pubsub.Client
	topicID string
}

var _ PubSub = (*cloudPB)(nil)

func NewCloudPubSub(client pubsub.Client, topicID string) *cloudPB {
	return &cloudPB{client: client, topicID: topicID}
}

func (p *cloudPB) Publish(ctx context.Context, msg []byte) error {
	if _, err := p.client.Topic(p.topicID).Publish(ctx, &pubsub.Message{Data: msg}).Get(ctx); err != nil {
		return ers.ErrInternal.New(err)
	}

	return nil
}

func (p *cloudPB) Subscribe(ctx context.Context, subFunc SubscribeFunc) error {
	id := fmt.Sprintf("%s-%s", p.topicID, uuid.New().String())
	sub, err := p.client.CreateSubscription(ctx, id, pubsub.SubscriptionConfig{Topic: p.client.Topic(p.topicID)})
	if err != nil {
		return ers.ErrInternal.New(err)
	}

	if err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		if err := subFunc(ctx, msg.Data); err != nil {
			return
		}
		msg.Ack()
	}); err != nil {
		return ers.W(err)
	}

	return nil
}

func (p *cloudPB) Close() error {
	if err := p.client.Close(); err != nil {
		return ers.W(err)
	}

	return nil
}

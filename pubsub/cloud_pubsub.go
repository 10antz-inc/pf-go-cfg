package pubsub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/10antz-inc/pf-go-ers"
	"github.com/google/uuid"
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
	topic, err := p.prepareTopic(ctx)
	if err != nil {
		return ers.W(err)
	}
	if _, err := topic.Publish(ctx, &pubsub.Message{Data: msg}).Get(ctx); err != nil {
		return ers.ErrInternal.New(err)
	}

	return nil
}

func (p *cloudPB) Subscribe(ctx context.Context, subFunc SubscribeFunc) error {
	topic, err := p.prepareTopic(ctx)
	if err != nil {
		return ers.W(err)
	}
	subID := fmt.Sprintf("%s-%s", p.topicID, uuid.NewString())
	sub, err := p.client.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{Topic: topic})
	if err != nil {
		return ers.W(err)
	}
	return ers.W(p.receive(ctx, sub, subFunc))
}

func (p *cloudPB) receive(ctx context.Context, sub *pubsub.Subscription, subFunc SubscribeFunc) error {
	var receiveError, funcError error
	receiveError = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		funcError = subFunc(ctx, msg.Data)
		if funcError != nil {
			return
		}
		msg.Ack()
	})
	if funcError != nil {
		return ers.W(funcError)
	}
	if receiveError != nil {
		return ers.W(receiveError)
	}
	return nil
}

func (p *cloudPB) Close() error {
	if err := p.client.Close(); err != nil {
		return ers.W(err)
	}

	return nil
}

func (p *cloudPB) prepareTopic(ctx context.Context) (*pubsub.Topic, error) {
	topic := p.client.Topic(p.topicID)
	exists, err := topic.Exists(ctx)
	if err != nil {
		return nil, ers.W(err)
	}
	if !exists {
		topic, err = p.client.CreateTopic(ctx, p.topicID)
		if err != nil {
			return nil, ers.W(err)
		}
	}
	return topic, nil
}

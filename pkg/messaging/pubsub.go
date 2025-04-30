package messaging

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"github.com/erdosam/go-clear-architecture/config"
	"github.com/erdosam/go-clear-architecture/pkg/logger"
)

type Pubsub interface {
	CreateTopics(topics []PubsubTopic)
	Publish(t PubsubTopic, d interface{}) error
	Subscribe(t PubsubTopic, s string, h handler) error
}

type PubsubTopic string

type handler func(context.Context, *pubsub.Message)

type googlePubsub struct {
	client *pubsub.Client
	log    logger.Interface
}

var _ Pubsub = &googlePubsub{}

func NewPubsub(cfg *config.Config, l logger.Interface) Pubsub {
	ctx := context.Background()
	c, err := pubsub.NewClient(ctx, cfg.Google.ProjectId)
	if err != nil {
		l.Fatal(err)
	}
	return &googlePubsub{client: c, log: l}
}

func (g *googlePubsub) CreateTopics(topics []PubsubTopic) {
	ctx := context.Background()
	for _, t := range topics {
		yes, err := g.client.Topic(string(t)).Exists(ctx)
		g.log.Info("creating topic %s", t)
		if err != nil {
			g.log.Debug(err)
		}
		if yes {
			continue
		}
		_, err = g.client.CreateTopic(ctx, string(t))
		if err != nil {
			g.log.Debug(err)
			continue
		}
		g.log.Info("%s topic created", t)
	}
}

func (g *googlePubsub) Publish(t PubsubTopic, d interface{}) error {
	tpc := g.client.Topic(string(t))
	defer tpc.Stop()

	data, _ := json.Marshal(d)
	msg := &pubsub.Message{Data: data}
	res := tpc.Publish(context.Background(), msg)
	g.log.Debug("Publish message to %s", string(t))
	if _, err := res.Get(context.Background()); err != nil {
		g.log.Error(err)
		return err
	}
	g.log.Info("Message published to %s: %s", string(t), data)
	return nil
}

func (g *googlePubsub) Subscribe(t PubsubTopic, s string, h handler) error {
	tpc := g.client.Topic(string(t))
	defer tpc.Stop()

	ctx := context.Background()
	sub := g.client.Subscription(s)
	yes, err := sub.Exists(ctx)
	if err != nil {
		g.log.Fatal(err)
	}
	if !yes {
		sub, err = g.client.CreateSubscription(ctx, s, pubsub.SubscriptionConfig{Topic: tpc})
		if err != nil {
			g.log.Fatal(err)
		}
		g.log.Info("%s subscription created", sub.ID())
	}
	return sub.Receive(ctx, h)
}

package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/erdosam/go-clear-architecture/internal/usecase"
	"github.com/erdosam/go-clear-architecture/pkg/logger"
	"github.com/erdosam/go-clear-architecture/pkg/messaging"
)

type SubscriptionHandler struct {
	log    logger.Interface
	pubsub messaging.Pubsub
}

func NewSubscriptionsHandler(l logger.Interface, ps messaging.Pubsub) *SubscriptionHandler {
	h := &SubscriptionHandler{
		log:    l,
		pubsub: ps,
	}
	return h
}

func (sh *SubscriptionHandler) Init() {
	sh.pubsub.CreateTopics([]messaging.PubsubTopic{
		usecase.PingTestTopic,
	})
	go sh.subscribePingTest()
}

func (sh *SubscriptionHandler) subscribePingTest() {
	//TODO must be in usecase
	t := usecase.PingTestTopic
	err := sh.pubsub.Subscribe(t, "subscription-"+string(t), func(_ context.Context, msg *pubsub.Message) {
		sh.log.Debug("Receive message: %s", msg.Data)
		msg.Ack()
	})
	if err != nil {
		sh.log.Error(err)
	}
}

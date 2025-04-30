package usecase

import "github.com/erdosam/go-clear-architecture/pkg/messaging"

type pingUsecase struct {
	pubsub messaging.Pubsub
}

const (
	PingTestTopic messaging.PubsubTopic = "ping_test"
	//TODO list down all topics here
)

var _ Ping = &pingUsecase{}

func NewPingUsecase(ps messaging.Pubsub) Ping {
	return &pingUsecase{ps}
}

func (p *pingUsecase) TestPubsub(msg string) error {
	return p.pubsub.Publish(PingTestTopic, msg)
}

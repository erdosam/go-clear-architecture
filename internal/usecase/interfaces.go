package usecase

import "github.com/erdosam/go-clear-architecture/internal/entity"

type Identity interface {
	GetToken() string
	GetId() string
}

// features
type (
	Ping interface {
		TestPubsub(msg string) error
	}
	User interface {
		GetUserFromId(id string, clientKey string) (entity.User, error)
	}
	//TODO define usecase you need (think like module or feature)
)

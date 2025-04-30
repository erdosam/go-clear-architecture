package webapi

import "github.com/erdosam/go-clear-architecture/internal/entity"

type (
	UserService interface {
		GetAccount() (account *entity.UserAccount, err error)
	}
)

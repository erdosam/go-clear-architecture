package webapi

import "github.com/arendi-project/ba-version-2/internal/entity"

type (
	UserService interface {
		GetAccount() (account *entity.UserAccount, err error)
	}
)

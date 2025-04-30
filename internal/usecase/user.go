package usecase

import (
	"github.com/erdosam/go-clear-architecture/internal/entity"
	repo "github.com/erdosam/go-clear-architecture/internal/usecase/dao"
	"github.com/erdosam/go-clear-architecture/pkg/logger"
)

var _ User = &userUseCase{}

type userUseCase struct {
	userDAO repo.UserDAO
	log     logger.Interface
}

func NewUserUseCase(l logger.Interface, dao repo.UserDAO) User {
	return &userUseCase{
		userDAO: dao,
		log:     l,
	}
}

func (u *userUseCase) GetUserFromId(id string, clientKey string) (entity.User, error) {
	user, err := u.userDAO.FindUserByUserId(id, clientKey)
	if err != nil {
		u.log.Error(err)
		return entity.User{}, err
	}
	return user, nil
}

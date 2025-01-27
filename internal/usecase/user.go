package usecase

import (
	"github.com/arendi-project/ba-version-2/internal/entity"
	repo "github.com/arendi-project/ba-version-2/internal/usecase/dao"
	"github.com/arendi-project/ba-version-2/pkg/logger"
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

func (u *userUseCase) GetUserByJunoId(id string, clientKey string) (entity.User, error) {
	user, err := u.userDAO.FindUserByJunoId(id, clientKey)
	if err != nil {
		u.log.Error(err)
		return entity.User{}, err
	}
	return user, nil
}

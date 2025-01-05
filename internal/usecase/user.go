package usecase

import (
	"github.com/arendi-project/ba-version-2/internal/entity"
	repo "github.com/arendi-project/ba-version-2/internal/usecase/dao"
	"github.com/arendi-project/ba-version-2/pkg/logger"
)

type userUseCase struct {
	userDAO repo.UserDAO
	log     logger.Interface
}

func NewUserUseCase(l logger.Interface, dao repo.UserDAO) *userUseCase {
	return &userUseCase{
		userDAO: dao,
		log:     l,
	}
}

func (u *userUseCase) GetUserById(id string) (entity.User, error) {
	return entity.User{Id: id}, nil
}

package dao

import (
	"errors"
	"github.com/arendi-project/ba-version-2/internal/entity"
	"github.com/arendi-project/ba-version-2/pkg/logger"
	"github.com/arendi-project/ba-version-2/pkg/postgres"
)

var _ UserDAO = &userDAO{}

type userDAO struct {
	*postgres.Postgres
	log logger.Interface
}

func NewUserDAO(l logger.Interface, pg *postgres.Postgres) UserDAO {
	return &userDAO{pg, l}
}

func (u *userDAO) FindUserById(id string) (entity.User, error) {
	if id == "" {
		return entity.User{}, errors.New("empty user")
	}
	u.log.Debug("Finding user", id)
	return entity.User{Id: id}, nil
}

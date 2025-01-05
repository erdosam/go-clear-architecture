package dao

import (
	"github.com/arendi-project/ba-version-2/internal/entity"
	"github.com/arendi-project/ba-version-2/pkg/logger"
	"github.com/arendi-project/ba-version-2/pkg/postgres"
)

type userDAO struct {
	*postgres.Postgres
	log logger.Interface
}

func NewUserDAO(l logger.Interface, pg *postgres.Postgres) *userDAO {
	return &userDAO{pg, l}
}

func (u *userDAO) FindUserById(id string) (*entity.User, error) {
	u.log.Debug("Finding user", id)
	return nil, nil
}

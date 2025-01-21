package dao

import (
	"errors"
	"fmt"
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

	var user entity.User
	u.log.Debug("Finding user with id %s", id)
	q := u.Rebind(`SELECT * FROM public.user WHERE id = ?`)
	if err := u.Get(&user, q, id); err != nil {
		return entity.User{}, fmt.Errorf("user %s not found", id)
	}
	return user, nil
}

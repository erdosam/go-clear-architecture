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

func (u *userDAO) FindUserByJunoId(id string, clientKey string) (entity.User, error) {
	if id == "" {
		return entity.User{}, errors.New("empty user")
	}

	var user entity.User
	u.log.Debug("Finding user with id %s and key %s", id, clientKey)
	q := u.Rebind(`SELECT
		u.id, u.display_name, u.status, u.mobile_code, u.phone_number
		FROM public.user u 
		    JOIN public.user_auth_juno j ON j.user_id = u.id WHERE j.account_id = ? AND j.client_key = ?`)
	if err := u.Get(&user, q, id, clientKey); err != nil {
		return entity.User{}, fmt.Errorf("user %s not found", id)
	}
	return user, nil
}

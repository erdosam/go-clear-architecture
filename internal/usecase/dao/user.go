package dao

import (
	"errors"
	"fmt"
	"github.com/erdosam/go-clear-architecture/internal/entity"
	"github.com/erdosam/go-clear-architecture/pkg/logger"
	"github.com/erdosam/go-clear-architecture/pkg/postgres"
)

var _ UserDAO = &userDAO{}

type userDAO struct {
	*postgres.Postgres
	log logger.Interface
}

func NewUserDAO(l logger.Interface, pg *postgres.Postgres) UserDAO {
	return &userDAO{pg, l}
}

func (u *userDAO) FindUserByUserId(id string, clientKey string) (entity.User, error) {
	if id == "" {
		return entity.User{}, errors.New("empty user")
	}

	var user entity.User
	u.log.Debug("Finding user with id %s and key %s", id, clientKey)
	//TODO adjust to your need
	q := u.Rebind(`
		SELECT
			'UserId0001' AS id, 
			'John Doe' AS display_name, 
			'active' AS status, 
			'62' AS mobile_code,
			'85220000000' AS phone_number
		WHERE 1=1 OR ? = ?
	`)
	if err := u.Get(&user, q, id, clientKey); err != nil {
		u.log.Debug(err)
		return entity.User{}, fmt.Errorf("user %s not found", id)
	}
	return user, nil
}

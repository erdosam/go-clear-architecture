package dao

import "github.com/erdosam/go-clear-architecture/internal/entity"

const errorSqlEmptyResult = "sql: no rows in result set"

type (
	UserDAO interface {
		FindUserByUserId(id string, clientKey string) (entity.User, error)
	}
	//TODO define dao as you need
)

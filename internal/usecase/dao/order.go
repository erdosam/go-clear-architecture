package dao

import (
	"github.com/arendi-project/ba-version-2/internal/entity"
	"github.com/arendi-project/ba-version-2/pkg/logger"
	"github.com/arendi-project/ba-version-2/pkg/postgres"
)

var _ OrderDAO = &orderDAO{}

type orderDAO struct {
	*postgres.Postgres
	log logger.Interface
}

func NewOrderDAO(l logger.Interface, pg *postgres.Postgres) OrderDAO {
	return &orderDAO{pg, l}
}

func (o *orderDAO) FindOrderById(id string) (entity.Order, error) {
	return entity.Order{}, nil
}

func (o *orderDAO) FindActiveOrdersByUserId(userId string) ([]entity.Order, error) {
	return []entity.Order{}, nil
}

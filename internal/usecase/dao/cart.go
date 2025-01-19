package dao

import (
	"github.com/arendi-project/ba-version-2/internal/entity"
	"github.com/arendi-project/ba-version-2/pkg/logger"
	"github.com/arendi-project/ba-version-2/pkg/postgres"
)

var _ CartingDAO = &cartingDAO{}

type cartingDAO struct {
	*postgres.Postgres
	log logger.Interface
}

func NewCartingDAO(log logger.Interface, pg *postgres.Postgres) CartingDAO {
	return &cartingDAO{pg, log}
}

func (dao *cartingDAO) FindItemsByCart(c entity.Cart) ([]entity.CartItem, error) {
	var rows []entity.CartItem
	q := dao.Rebind(`SELECT * FROM public.cart_item WHERE user_id = ?`)
	err := dao.Select(&rows, q, c.UserId)
	if err != nil {
		dao.log.Error(err)
		return nil, err
	}
	dao.log.Debug("Items count %d", len(rows))
	return rows, nil
}

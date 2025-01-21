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
		panic(err.(any))
	}
	dao.log.Debug("User %s's items count %d", c.UserId, len(rows))
	return rows, nil
}

func (dao *cartingDAO) FindOneItem(arg ...interface{}) (entity.CartItem, error) {
	var item entity.CartItem
	query, args := buildFindQuery(`SELECT * FROM public.cart_item`, arg...)
	q := dao.Rebind(query)
	if err := dao.Get(&item, q, args...); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return entity.CartItem{}, err
		}
		panic(err.(any))
	}
	return item, nil
}

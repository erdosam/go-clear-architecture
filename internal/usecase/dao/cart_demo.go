//go:build demo

package dao

import (
	"github.com/arendi-project/ba-version-2/internal/entity"
	"github.com/arendi-project/ba-version-2/pkg/postgres"
)
import "github.com/arendi-project/ba-version-2/pkg/logger"

var _ CartingDAO = &cartingDemoDAO{}

type cartingDemoDAO struct {
	*postgres.Postgres
	log logger.Interface
}

func NewCartingDAO(log logger.Interface, pg *postgres.Postgres) CartingDAO {
	return &cartingDemoDAO{pg, log}
}

func (dao *cartingDemoDAO) FindItemsByCart(c entity.Cart) ([]entity.CartItem, error) {
	dao.log.Info("Demo items count %d", 0)
	return []entity.CartItem{}, nil
}

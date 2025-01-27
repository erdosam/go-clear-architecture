package dao

import "github.com/arendi-project/ba-version-2/internal/entity"

const errorSqlEmptyResult = "sql: no rows in result set"

type (
	UserDAO interface {
		FindUserByJunoId(id string, clientKey string) (entity.User, error)
	}
	CartingDAO interface {
		FindItemsByCart(c entity.Cart) ([]entity.CartItem, error)
		FindOneItem(arg ...interface{}) (entity.CartItem, error)
	}
	OrderDAO interface {
		FindOrderById(id string) (entity.Order, error)
		FindActiveOrdersByUserId(userId string) ([]entity.Order, error)
	}
	TrashCategoryDAO interface {
		FindCategories(arg ...interface{}) ([]entity.TrashCategory, error)
	}
)

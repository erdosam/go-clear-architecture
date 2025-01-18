package dao

import "github.com/arendi-project/ba-version-2/internal/entity"

type (
	UserDAO interface {
		FindUserById(id string) (entity.User, error)
	}
	CartingDAO interface {
		FindItemsByCart(c entity.Cart) ([]entity.CartItem, error)
	}
	OrderDAO interface {
		FindOrderById(id string) (entity.Order, error)
		FindActiveOrdersByUserId(userId string) ([]entity.Order, error)
	}
)

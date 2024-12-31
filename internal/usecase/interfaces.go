package usecase

import "github.com/arendi-project/ba-version-2/internal/entity"

// features
type (
	User interface {
	}
	Carting interface {
		GetItems(c entity.Cart) ([]entity.CartItem, error)
		AddItemToCart(i entity.CartItem, c entity.Cart) error
	}
	Order interface {
		CreateMultipleItemsOrder(i []entity.CartItem) (entity.Order, error)
		GetOrderById(id string) (entity.Order, error)
	}
)

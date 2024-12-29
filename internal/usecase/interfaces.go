package usecase

import "github.com/arendi-project/ba-version-2/internal/entity"

// form models
type (
	User interface {
	}
	Carting interface {
		GetItems(c entity.Cart) ([]entity.CartItem, error)
		AddItemToCart(i entity.CartItem, c entity.Cart) error
	}
)

// DAOs
type (
	UserDAO interface {
		GetUserById(id string) (entity.User, error)
	}
	CartingDAO interface {
		GetItemsByCart(c entity.Cart) ([]entity.CartItem, error)
	}
)

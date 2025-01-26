package usecase

import "github.com/arendi-project/ba-version-2/internal/entity"

type Identity interface {
	GetToken() string
	GetId() string
}

// features
type (
	User interface {
		GetUserById(id string) (entity.User, error)
	}
	Category interface {
		GetAvailableCategories() ([]entity.TrashCategory, error)
	}
	Carting interface {
		GetCart(u entity.User) (entity.Cart, error)
		GetItems(c entity.Cart) ([]entity.CartItem, error)
		GetItem(c entity.Cart, id string) (entity.CartItem, error)
		AddItemToCart(req entity.AddItemToCartForm) error
		EditCartItem(req entity.EditCartItemForm) error
	}
	Order interface {
		CreateMultipleItemsOrder(i []entity.CartItem) (entity.Order, error)
		GetOrderById(id string) (entity.Order, error)
	}
)

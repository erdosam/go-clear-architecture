package entity

import "github.com/go-playground/validator/v10"

var _ baseFormRequest = &AddItemToCartForm{}

type (
	Cart struct {
		UserId string `json:"id"`
		Timestamp
	}
	CartItem struct {
		Id         string `db:"id" json:"id"`
		UserId     string `db:"user_id" json:"user_id"`
		CategoryId string `db:"category_id" json:"category_id"`
		Quantity   string `db:"quantity" json:"quantity"`
		Blamable
	}
	AddItemToCartForm struct {
		CategoryId string `json:"category_id" validate:"required" error-required:"Missing category_id field"`
		Quantity   int    `json:"quantity" validate:"required,gte=1" error-required:"Missing quantity field" error-gte:"Quantity must be greater than 1"`
	}
	EditCartItemForm struct {
		ItemId   string `json:"id" validate:"required" error-required:"Missing id field for item id"`
		Quantity int    `json:"quantity" validate:"required,gte=1" error-required:"Missing quantity field" error-gte:"Quantity must be at least 1"`
	}
	CartItemResponse struct {
		Items []CartItem `json:"items"`
	}
)

func (o *AddItemToCartForm) Validate(v *validator.Validate) error {
	return validateFunc[AddItemToCartForm](*o, v)
}

func (o *EditCartItemForm) Validate(v *validator.Validate) error {
	return validateFunc[EditCartItemForm](*o, v)
}

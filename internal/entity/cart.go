package entity

import "github.com/go-playground/validator/v10"

var _ baseFormRequest = &AddItemToCartRequest{}

type (
	Cart struct {
		UserId string `json:"id"`
		Timestamp
	}
	CartItem struct {
		Id         string `db:"id" json:"id"`
		UserId     string `db:"user_id" json:"user_id"`
		CategoryId string `json:"category_id"`
		Quantity   string `json:"quantity"`
		Blamable
	}
	AddItemToCartRequest struct {
		CategoryId string `json:"category_id" validate:"required" error-required:"Missing category_id field"`
		Quantity   int    `json:"quantity" validate:"required,gte=1" error-required:"Missing quantity field" error-gte:"Quantity must be greater than 1"`
	}
	CartItemResponse struct {
		Items []CartItem `json:"items"`
	}
)

func (m *AddItemToCartRequest) Validate(v *validator.Validate) error {
	return validateFunc[AddItemToCartRequest](*m, v)
}

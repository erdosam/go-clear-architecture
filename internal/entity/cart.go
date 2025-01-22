package entity

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
		baseFormRequest[AddItemToCartRequest]
		ItemId string `json:"item_id" validate:"required" error-required:"Item id is required"`
		CartId string `json:"cart_id" validate:"required" error-required:"Cart id is required"`
		Amount int    `json:"amount" validate:"required,gte=2000" error-required:"Field amount is required" error-gte:"Amount must be greater than 2000"`
	}
	CartItemResponse struct {
		Items []CartItem `json:"items"`
	}
)

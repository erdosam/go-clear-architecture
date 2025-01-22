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
		ItemId string `json:"item_id" validate:"required"`
		CartId string `json:"cart_id" validate:"required"`
	}
	CartItemResponse struct {
		Items []CartItem `json:"items"`
	}
)

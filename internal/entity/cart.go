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
	CartItemResponse struct {
		Items []CartItem `json:"items"`
	}
)

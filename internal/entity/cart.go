package entity

type (
	User struct {
		Id string `db:"id" json:"id"`
	}
	UserAddress struct {
		Id        string `db:"id" json:"id"`
		UserId    string `db:"user_id" json:"user_id"`
		Details   string `db:"detail" json:"details"`
		GeoHash   string `db:"geohash"`
		CreatedAt string `db:"created_at" json:"created_at"`
		DeletedAt string `db:"deleted_at" json:"deleted_at"`
	}
)

type (
	Cart struct {
		UserId string `json:"id"`
	}
	CartItem struct {
		Id         int    `db:"id" json:"id"`
		UserId     string `db:"user_id" json:"user_id"`
		CategoryId string `json:"category_id"`
		Quantity   string `json:"quantity"`
		CreatedAt  string `json:"created_at"`
		CreatedBy  string `json:"created_by"`
	}
)

package entity

type (
	User struct {
		Id          string `db:"id" json:"id"`
		Name        string `db:"display_name" json:"name"`
		MobileCode  string `db:"mobile_code"`
		PhoneNumber string `db:"phone_number"`
		Status      string `db:"status"`
	}
	UserAccount struct {
		UserId string `db:"user_id"`
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

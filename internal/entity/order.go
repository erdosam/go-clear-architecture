package entity

type Order struct {
	Id string `db:"id" json:"id"`
	Timestamp
}

type OrderResponse struct {
}

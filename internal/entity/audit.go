package entity

type Timestamp struct {
	CreatedAt string `db:"created_at" json:"created_at"`
	CreatedBy string `db:"created_by" json:"created_by"`
}

type Blamable struct {
	Timestamp
	UpdatedAt string `db:"updated_at" json:"updated_at"`
	UpdatedBy string `db:"updated_by" json:"updated_by"`
}

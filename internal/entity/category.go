package entity

type TrashCategoryGroup string

type TrashCategory struct {
	Id               string             `db:"id" json:"id"`
	Name             string             `db:"name" json:"name"`
	ParentCategoryId *string            `db:"parent_category_id" json:"parent_category_id"`
	Group            TrashCategoryGroup `db:"group" json:"group"`
	Status           string             `db:"status" json:"status"`
	Description      string             `db:"description" json:"description"`
	Image            ImageUrl           `db:"image" json:"image"`
	Timestamp
}

type TrashCategoryResponse struct {
	*TrashCategory
}

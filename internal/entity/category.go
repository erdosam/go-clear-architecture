package entity

type TrashCategoryGroup string

type TrashCategory struct {
	Id               string             `db:"id"`
	Name             string             `db:"name"`
	ParentCategoryId string             `db:"parent_category_id"`
	Group            TrashCategoryGroup `db:"group"`
	Status           string             `db:"status"`
	Timestamp
}

type TrashCategoryDetail struct {
	Id             string `db:"id"`
	CategoryId     string `db:"category_id"`
	Description    string `db:"description"`
	ImageUrlMedium string `db:"image_url_md"`
	ImageUrlSmall  string `db:"image_url_sm"`
	ImageUrlLarge  string `db:"image_url_lg"`
	Timestamp
}

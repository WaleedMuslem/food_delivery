package model

type Category struct {
	Category_id    int    `json:"category_id"`
	Category_name  string `json:"category_name"`
	Category_image string `json:"category_image"`
	ProductCount   int    `json:"category_count"`
}

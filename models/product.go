package models

type Product struct {
	ProductId      uint    `gorm:"primaryKey" json:"product_id"`
	Name           string  `json:"name"`
	Price          float64 `json:"price"`
	Description    string  `json:"description"`
	CategoryName   string  `json:"category_name"`
	ImageUrl       string  `json:"image_url"`
	ImageLocalPath string  `json:"image_local_path"`
	Stock          string  `json:"stock"`
}

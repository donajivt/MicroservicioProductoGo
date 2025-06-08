package dtos

import "mime/multipart"

type ProductDto struct {
	ProductId      uint                  `json:"product_id"`
	Name           string                `json:"name"`
	Price          float64               `json:"price"`
	Description    string                `json:"description"`
	CategoryName   string                `json:"category_name"`
	ImageUrl       string                `json:"image_url"`
	ImageLocalPath string                `json:"image_local_path"`
	Image          *multipart.FileHeader `json:"-"`
	Stock          string                `json:"stock"`
}

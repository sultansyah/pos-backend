package product

import (
	"post-backend/internal/category"
	"time"
)

type Product struct {
	Id            int               `json:"id"`
	CategoryId    int               `json:"category_id"`
	Category      category.Category `json:"category"`
	ProductImages []ProductImages   `json:"product_images"`
	Name          string            `json:"name"`
	Slug          string            `json:"slug"`
	Description   string            `json:"description"`
	Stock         int               `json:"stock"`
	Price         float64           `json:"price"`
	Status        string            `json:"status"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type ProductImages struct {
	Id        int       `json:"id"`
	ProductId int       `json:"product_id"`
	ImageUrl  string    `json:"image_url"`
	IsLogo    bool      `json:"is_logo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

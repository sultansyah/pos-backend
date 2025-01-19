package product

type CreateProductInput struct {
	CategoryId  int     `form:"category_id" binding:"required"`
	Name        string  `form:"name" binding:"required"`
	Description string  `form:"description" binding:"required"`
	Stock       int     `form:"stock" binding:"required"`
	Price       float64 `form:"price" binding:"required"`
	Status      string  `form:"status" binding:"required"`
}

type UpdateProductInput struct {
	CategoryId  int     `json:"category_id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Stock       int     `json:"stock" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Status      string  `json:"status" binding:"required"`
}

type GetProductInput struct {
	Id int `uri:"id" binding:"required"`
}

type SlugProductInput struct {
	Slug string `uri:"slug" binding:"required"`
}

type GetProductImageInput struct {
	Id int `uri:"imageId" binding:"required"`
}

type UpdateStockProductInput struct {
	Qty  int    `json:"qty" binding:"required"`
	Type string `json:"type" binding:"required"`
}

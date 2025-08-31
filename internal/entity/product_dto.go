package entity

// ProductCreate - DTO для запроса создания товара (POST /products).
type ProductCreate struct {
	Name       string  `json:"name"        validate:"required,min=3"`
	Price      float64 `json:"price"       validate:"required,gt=0"`
	CategoryID *int    `json:"category_id" validate:"omitempty,gt=0"`
}

// ProductUpdate - DTO для частичного обновления (PATCH /products/:id).
type ProductUpdate struct {
	ID         int      `json:"id"          validate:"required"`
	Name       *string  `json:"name"        validate:"omitempty,min=3"`
	Price      *float64 `json:"price"       validate:"omitempty,gt=0"`
	CategoryID *int     `json:"category_id" validate:"omitempty,gt=0"`
}

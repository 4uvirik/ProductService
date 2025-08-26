package entity

// Product - основная сущность товара, как в БД
type Product struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	CategoryID *int    `json:"category_id"`
}

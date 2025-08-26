package entity

// Category - сущность категорий товара, как в БД
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name" validate:"required,min=3"`
}

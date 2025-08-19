package entity

type Product struct {
	ID         int
	Name       string
	CategoryID *int
	Price      float64
}

type Category struct {
	ID   int
	Name string
}

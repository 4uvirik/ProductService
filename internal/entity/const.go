package entity

// ------ Константы для реализации методов ProductOperations ------

const QueryProductCreate = `
INSERT INTO products (name, price, category_id)
VALUES ($1, $2, $3)
RETURNING id;
`
const QueryProductGetAll = `
SELECT id, name, price, category_id
FROM products
ORDER BY id;
`
const QueryProductGetByID = `
SELECT id, name, price, category_id
FROM products
WHERE id = $1;
`
const QueryProductUpdate = `
UPDATE products
SET name = $1, price = $2, category_id = $3
WHERE id = $4;
`
const QueryProductUpdatePrice = `
UPDATE products 
SET price=$1 
WHERE id=$2
`
const QueryProductDelete = `
DELETE FROM products 
WHERE id = $1
`

// ------ Константы для реализации методов CategoryOperations ------

const QueryCategoryCreate = `
INSERT INTO category (name)
VALUES ($1)
RETURNING id;
`
const QueryCategoryUpdate = `
UPDATE category
SET name=$1 
WHERE id=$2
`
const QueryCategoryDelete = `
DELETE FROM category 
WHERE id = $1
`

// ------ Константы для URL ------

const ProductURL = "/product"
const CategoryURL = "/category"
const IDParam = "/:id"

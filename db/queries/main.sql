-- name: CreateUser :one
INSERT INTO users (
  username,
  email,
  hashed_password,
  full_name
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: CreateCategory :one
INSERT INTO categories (
  name
) VALUES (
  $1
) RETURNING *;

-- name: CreateProduct :one
INSERT INTO products (
  name,
  description,
  price,
  stock,
  category_id
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: CreateOrder :one
INSERT INTO orders (
  user_id,
  total_price,
  status
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (
  order_id,
  product_id,
  quantity,
  price
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: UpdateOrder :one
UPDATE orders
SET total_price = $2, status = $3
WHERE id = $1
RETURNING *;

-- name: UpdateProductStock :one
UPDATE products
SET stock = stock - $2
WHERE id = $1
RETURNING *;

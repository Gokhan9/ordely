package domain

import (
	"context"
	"time"

	"github.com/gokhan/orderly/internal/repository/db"
)

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int32   `json:"stock" validate:"required,gte=0"`
	CategoryID  int64   `json:"category_id" validate:"required"`
}

type ProductResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	Stock       int32     `json:"stock"`
	CategoryID  int64     `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type ListProductsRequest struct {
	PageID   int32 `query:"page_id" validate:"required,min=1"`
	PageSize int32 `query:"page_size" validate:"required,min=5,max=20"`
}

type ProductUseCase interface {
	CreateProduct(ctx context.Context, req CreateProductRequest) (ProductResponse, error)
	GetProduct(ctx context.Context, id int64) (ProductResponse, error)
	ListProducts(ctx context.Context, req ListProductsRequest) ([]ProductResponse, error)
}

func NewProductResponse(product db.Product) ProductResponse {
	return ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price.String(),
		Stock:       product.Stock,
		CategoryID:  product.CategoryID,
		CreatedAt:   product.CreatedAt,
	}
}

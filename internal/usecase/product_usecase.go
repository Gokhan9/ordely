package usecase

import (
	"context"

	"github.com/gokhan/orderly/internal/domain"
	"github.com/gokhan/orderly/internal/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type productUseCase struct {
	store db.Store
}

func NewProductUseCase(store db.Store) domain.ProductUseCase {
	return &productUseCase{
		store: store,
	}
}

func (u *productUseCase) CreateProduct(ctx context.Context, req domain.CreateProductRequest) (domain.ProductResponse, error) {
	var price pgtype.Numeric
	price.Scan(req.Price)

	arg := db.CreateProductParams{
		Name:        req.Name,
		Description: req.Description,
		Price:       price,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
	}

	product, err := u.store.CreateProduct(ctx, arg)
	if err != nil {
		return domain.ProductResponse{}, err
	}

	return domain.NewProductResponse(product), nil
}

func (u *productUseCase) GetProduct(ctx context.Context, id int64) (domain.ProductResponse, error) {
	product, err := u.store.GetProduct(ctx, id)
	if err != nil {
		return domain.ProductResponse{}, err
	}

	return domain.NewProductResponse(product), nil
}

func (u *productUseCase) ListProducts(ctx context.Context, req domain.ListProductsRequest) ([]domain.ProductResponse, error) {
	arg := db.ListProductsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	products, err := u.store.ListProducts(ctx, arg)
	if err != nil {
		return nil, err
	}

	var resp []domain.ProductResponse
	for _, p := range products {
		resp = append(resp, domain.NewProductResponse(p))
	}

	return resp, nil
}

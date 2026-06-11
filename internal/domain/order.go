package domain

import (
	"context"

	"github.com/gokhan/orderly/internal/repository/db"
)

type OrderItemRequest struct {
	ProductID int64 `json:"product_id" validate:"required"`
	Quantity  int32 `json:"quantity" validate:"required,gt=0"`
}

type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items" validate:"required,dive"`
}

type OrderResponse struct {
	ID         int64   `json:"id"`
	UserID     int64   `json:"user_id"`
	TotalPrice string  `json:"total_price"`
	Status     string  `json:"status"`
}

type OrderUseCase interface {
	CreateOrder(ctx context.Context, userID int64, req CreateOrderRequest) (OrderResponse, error)
}

func NewOrderResponse(order db.Order) OrderResponse {
	return OrderResponse{
		ID:         order.ID,
		UserID:     order.UserID,
		TotalPrice: order.TotalPrice.String(),
		Status:     order.Status,
	}
}

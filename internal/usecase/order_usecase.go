package usecase

import (
	"context"
	"fmt"

	"github.com/gokhan/orderly/internal/domain"
	"github.com/gokhan/orderly/internal/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type orderUseCase struct {
	store db.Store
}

func NewOrderUseCase(store db.Store) domain.OrderUseCase {
	return &orderUseCase{
		store: store,
	}
}

func (u *orderUseCase) CreateOrder(ctx context.Context, userID int64, req domain.CreateOrderRequest) (domain.OrderResponse, error) {
	var orderResponse domain.OrderResponse

	err := u.store.ExecTx(ctx, func(q *db.Queries) error {
		var totalAmount float64

		// 1. Create Order Placeholder
		var totalPrice pgtype.Numeric
		totalPrice.Scan(0)
		order, err := q.CreateOrder(ctx, db.CreateOrderParams{
			UserID:     userID,
			TotalPrice: totalPrice,
			Status:     "pending",
		})
		if err != nil {
			return err
		}

		// 2. Process Items
		for _, itemReq := range req.Items {
			product, err := q.GetProduct(ctx, itemReq.ProductID)
			if err != nil {
				return fmt.Errorf("product %d not found", itemReq.ProductID)
			}

			if product.Stock < itemReq.Quantity {
				return fmt.Errorf("insufficient stock for product %d", itemReq.ProductID)
			}

			// Add order item
			itemPrice := product.Price
			_, err = q.CreateOrderItem(ctx, db.CreateOrderItemParams{
				Order_id:  order.ID,
				Product_id: product.ID,
				Quantity:  itemReq.Quantity,
				Price:     itemPrice,
			})
			if err != nil {
				return err
			}

			// Update Stock
			_, err = q.UpdateProductStock(ctx, db.UpdateProductStockParams{
				ID:    product.ID,
				Stock: itemReq.Quantity,
			})
			if err != nil {
				return err
			}

			// Update total amount
			f, _ := product.Price.Float64Value()
			totalAmount += f.Float64 * float64(itemReq.Quantity)
		}

		// 3. Update Order Total Price
		var finalTotalPrice pgtype.Numeric
		finalTotalPrice.Scan(fmt.Sprintf("%.2f", totalAmount))
		
		updatedOrder, err := q.UpdateOrder(ctx, db.UpdateOrderParams{
			ID:         order.ID,
			TotalPrice: finalTotalPrice,
			Status:     "completed",
		})
		if err != nil {
			return err
		}

		orderResponse = domain.NewOrderResponse(updatedOrder)
		return nil
	})

	return orderResponse, err
}

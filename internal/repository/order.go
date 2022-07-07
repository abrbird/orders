package repository

import (
	"context"
	"github.com/abrbird/orders/internal/models"
)

type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) error
	Retrieve(ctx context.Context, orderId int64) models.OrderRetrieve
	Update(ctx context.Context, order *models.Order) error

	CreateItem(ctx context.Context, orderItem *models.OrderItem) error
	RetrieveItem(ctx context.Context, orderItemId int64) models.OrderItemRetrieve
	RetrieveItems(ctx context.Context, orderId int64) models.OrderItemRetrieve
	UpdateItem(ctx context.Context, orderItem *models.OrderItem) error
}

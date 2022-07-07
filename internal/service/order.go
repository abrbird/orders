package service

import (
	"context"
	"gitlab.ozon.dev/zBlur/homework-3/orders/internal/models"
	"gitlab.ozon.dev/zBlur/homework-3/orders/internal/repository"
)

type OrderService interface {
	Create(ctx context.Context, repository repository.OrderRepository, order *models.Order) error
	Retrieve(ctx context.Context, repository repository.OrderRepository, orderId int64) models.OrderRetrieve
	Update(ctx context.Context, repository repository.OrderRepository, order *models.Order) error

	CreateItem(ctx context.Context, repository repository.OrderRepository, orderItem *models.OrderItem) error
	RetrieveItem(ctx context.Context, repository repository.OrderRepository, orderItemId int64) models.OrderItemRetrieve
	RetrieveItems(ctx context.Context, repository repository.OrderRepository, orderId int64) models.OrderItemsRetrieve
	UpdateItem(ctx context.Context, repository repository.OrderRepository, orderItem *models.OrderItem) error

	MarkOrderIssued(ctx context.Context, repository repository.OrderRepository, orderId int64) error
}

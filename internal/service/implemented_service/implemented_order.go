package implemented_service

import (
	"context"
	"gitlab.ozon.dev/zBlur/homework-3/orders/internal/models"
	"gitlab.ozon.dev/zBlur/homework-3/orders/internal/repository"
)

type OrderService struct{}

func (o OrderService) Create(ctx context.Context, repository repository.OrderRepository, order *models.Order) error {
	//TODO implement me
	panic("implement me")
}

func (o OrderService) Retrieve(ctx context.Context, repository repository.OrderRepository, orderId int64) models.OrderRetrieve {
	return repository.Retrieve(ctx, orderId)
}

func (o OrderService) Update(ctx context.Context, repository repository.OrderRepository, order *models.Order) error {
	return repository.Update(ctx, order)
}

func (o OrderService) CreateItem(ctx context.Context, repository repository.OrderRepository, orderItem *models.OrderItem) error {
	//TODO implement me
	panic("implement me")
}

func (o OrderService) RetrieveItem(ctx context.Context, repository repository.OrderRepository, orderItemId int64) models.OrderItemRetrieve {
	//TODO implement me
	panic("implement me")
}

func (o OrderService) RetrieveItems(ctx context.Context, repository repository.OrderRepository, orderId int64) models.OrderItemsRetrieve {
	//TODO implement me
	panic("implement me")
}

func (o OrderService) UpdateItem(ctx context.Context, repository repository.OrderRepository, orderItem *models.OrderItem) error {
	//TODO implement me
	panic("implement me")
}

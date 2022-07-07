package implemented_service

import (
	"context"
	"github.com/abrbird/orders/internal/cache/redis_cache"
	"github.com/abrbird/orders/internal/models"
	"github.com/abrbird/orders/internal/repository"
	"github.com/pkg/errors"
)

type OrderService struct {
	service *Service
}

func (o OrderService) Retrieve(ctx context.Context, repository repository.OrderRepository, orderId int64) models.OrderRetrieve {
	retrieved := o.service.cache.Order().Get(ctx, orderId)
	if errors.Is(retrieved.Error, redis_cache.Nil) {
		retrieved = repository.Retrieve(ctx, orderId)

		if retrieved.Error != nil {
			return retrieved
		}
	}

	if err := o.service.cache.Order().Set(ctx, *retrieved.Order); err != nil {
		retrieved.Order = nil
		retrieved.Error = err
		return retrieved
	}

	return retrieved
}

func (o OrderService) Update(ctx context.Context, repository repository.OrderRepository, order *models.Order) error {
	err := repository.Update(ctx, order)
	if err != nil {
		return err
	}

	if err = o.service.cache.Order().Set(ctx, *order); err != nil {
		return err
	}

	return nil
}

func (o OrderService) MarkOrderIssued(ctx context.Context, repository repository.OrderRepository, orderId int64) error {
	orderRetrieved := o.Retrieve(
		ctx,
		repository,
		orderId,
	)

	if orderRetrieved.Error != nil {
		//return models.NewRetryError(orderRetrieved.Error)
		return models.NewRetryError(nil)
	}

	orderRetrieved.Order.Status = models.Issued
	err := o.Update(
		ctx,
		repository,
		orderRetrieved.Order,
	)
	if err != nil {
		//return models.NewRetryError(orderRetrieved.Error)
		return models.NewRetryError(nil)
	}

	return nil
}

func (o OrderService) Create(ctx context.Context, repository repository.OrderRepository, order *models.Order) error {
	//TODO implement me
	panic("implement me")
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

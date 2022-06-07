package sql_repository

import (
	"context"
	"gitlab.ozon.dev/zBlur/homework-3/orders/internal/models"
)

type SQLOrderRepository struct {
	store *SQLRepository
}

func (S SQLOrderRepository) Create(ctx context.Context, order *models.Order) error {
	//TODO implement me
	panic("implement me")
}

func (S SQLOrderRepository) Retrieve(ctx context.Context, orderId int64) models.OrderRetrieve {
	const query = `
		SELECT 
    		id,
			status
		FROM orders_order
		WHERE id = $1
	`

	order := &models.Order{}
	if err := S.store.dbConnectionPool.QueryRow(
		ctx,
		query,
		orderId,
	).Scan(
		&order.Id,
		&order.Status,
	); err != nil {
		return models.OrderRetrieve{Order: nil, Error: models.NotFoundError}
	}
	return models.OrderRetrieve{Order: order, Error: nil}
}

func (S SQLOrderRepository) Update(ctx context.Context, order *models.Order) error {
	const query = `
		UPDATE orders_order
		SET (status) = ($2)
		WHERE id = $1 
	`

	err := S.store.dbConnectionPool.QueryRow(
		ctx,
		query,
		order.Id,
		order.Status,
	)
	if err != nil {
		return models.NotFoundError
	}
	return nil
}

func (S SQLOrderRepository) CreateItem(ctx context.Context, orderItem *models.OrderItem) error {
	//TODO implement me
	panic("implement me")
}

func (S SQLOrderRepository) RetrieveItem(ctx context.Context, orderItemId int64) models.OrderItemRetrieve {
	//TODO implement me
	panic("implement me")
}

func (S SQLOrderRepository) RetrieveItems(ctx context.Context, orderId int64) models.OrderItemRetrieve {
	//TODO implement me
	panic("implement me")
}

func (S SQLOrderRepository) UpdateItem(ctx context.Context, orderItem *models.OrderItem) error {
	//TODO implement me
	panic("implement me")
}

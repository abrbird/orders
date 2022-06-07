package sql_repository

import (
	"context"
	"gitlab.ozon.dev/zBlur/homework-3/orders/internal/models"
)

type SQLOrderRepository struct {
	store *SQLRepository
}

//type SQLUser struct {
//	Id        int64
//	UserName  sql.NullString
//	FirstName sql.NullString
//	LastName  sql.NullString
//}

func (S SQLOrderRepository) Create(ctx context.Context, order *models.Order) error {
	//TODO implement me
	panic("implement me")
}

func (S SQLOrderRepository) Retrieve(ctx context.Context, orderId int64) models.OrderRetrieve {
	//TODO implement me
	panic("implement me")
}

func (S SQLOrderRepository) Update(ctx context.Context, order *models.Order) error {
	//TODO implement me
	panic("implement me")
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

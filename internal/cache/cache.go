package cache

import (
	"context"
	"gitlab.ozon.dev/zBlur/homework-3/orders/internal/models"
)

type Cache interface {
	Order() OrderCache
}

type OrderCache interface {
	Get(ctx context.Context, id int64) models.OrderRetrieve
	Set(ctx context.Context, order models.Order) error
}

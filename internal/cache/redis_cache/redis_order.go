package redis_cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/abrbird/orders/internal/models"
	"time"
)

type RedisOrderCache struct {
	Prefix     string
	redisCache *RedisCache
	expiration time.Duration
}

func (r RedisOrderCache) getIDKey(id int64) string {
	return fmt.Sprintf("%s_%v", r.Prefix, id)
}

func (r RedisOrderCache) Get(ctx context.Context, id int64) models.OrderRetrieve {
	//
	//	span, ctx := opentracing.StartSpanFromContext(ctx, "cache")
	//	defer span.Finish()
	//

	item, err := r.redisCache.client.Get(ctx, r.getIDKey(id)).Result()
	if err != nil {
		return models.OrderRetrieve{Order: nil, Error: err}
	}

	var record models.Order
	if err = json.Unmarshal([]byte(item), &record); err != nil {
		return models.OrderRetrieve{Order: nil, Error: err}
	}

	return models.OrderRetrieve{Order: &record, Error: nil}
}

func (r RedisOrderCache) Set(ctx context.Context, order models.Order) error {
	//
	//	span, ctx := opentracing.StartSpanFromContext(ctx, "cache")
	//	defer span.Finish()
	//

	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	return r.redisCache.client.Set(ctx, r.getIDKey(order.Id), data, r.expiration).Err()
}

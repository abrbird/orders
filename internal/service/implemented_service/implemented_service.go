package implemented_service

import (
	"gitlab.ozon.dev/zBlur/homework-3/orders/internal/cache"
	"gitlab.ozon.dev/zBlur/homework-3/orders/internal/service"
)

type Service struct {
	cache        cache.Cache
	orderService *OrderService
}

func New(cache cache.Cache) *Service {
	srvc := &Service{
		cache: cache,
	}
	srvc.orderService = &OrderService{srvc}

	return srvc
}

func (s *Service) Order() service.OrderService {
	return s.orderService
}

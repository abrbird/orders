package implemented_service

import "gitlab.ozon.dev/zBlur/homework-3/orders/internal/service"

type Service struct {
	orderService *OrderService
}

func New() *Service {
	return &Service{
		orderService: &OrderService{},
	}
}

func (s *Service) Order() service.OrderService {
	return s.orderService
}

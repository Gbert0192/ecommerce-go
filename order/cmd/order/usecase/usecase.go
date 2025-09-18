package usecase

import "order/order/cmd/order/service"

type OrderUseCase struct {
	OrderService *service.OrderService
}

func NewOrderUseCase(orderService *service.OrderService) *OrderUseCase {
	return &OrderUseCase{
		OrderService: orderService,
	}
}

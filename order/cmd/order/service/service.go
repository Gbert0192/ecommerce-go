package service

import "order/order/cmd/order/repository"

type OrderService struct {
	OrderRepository *repository.OrderRepository
}

func NewOrderService(orderRepo *repository.OrderRepository) *OrderService {
	return &OrderService{
		OrderRepository: orderRepo,
	}
}

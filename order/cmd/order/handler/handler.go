package handler

import "order/order/cmd/order/usecase"

type OrderHandler struct {
	OrderUsecase *usecase.OrderUseCase
}

func NewOrderHandler(orderUseCase *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		OrderUsecase: orderUseCase,
	}
}

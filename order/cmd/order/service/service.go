package service

import (
	"context"
	"errors"
	"order/order/cmd/order/repository"
	"order/order/models"

	"gorm.io/gorm"
)

type OrderService struct {
	OrderRepository *repository.OrderRepository
}

func NewOrderService(orderRepo *repository.OrderRepository) *OrderService {
	return &OrderService{
		OrderRepository: orderRepo,
	}
}

func (s *OrderService) CheckIdempotency(ctx context.Context, token string) (bool, error) {
	isExists, err := s.OrderRepository.CheckIdempotency(ctx, token)
	if err != nil {
		return false, err
	}
	return isExists, err
}

func (s *OrderService) SaveIdempotency(ctx context.Context, token string) error {
	err := s.OrderRepository.SaveIdempotency(ctx, token)
	if err != nil {
		return err
	}
	return nil
}

func (s *OrderService) SaveOrderAndOrderDetail(ctx context.Context, order *models.Order, orderDetail *models.OrderDetail) (int64, error) {
	var orderID int64

	err := s.OrderRepository.WithTransaction(ctx, func(tx *gorm.DB) error {
		err := s.OrderRepository.InsertOrderDetailTx(ctx, tx, orderDetail)
		if err != nil {
			return err
		}
		order.OrderDetailID = orderDetail.ID
		err = s.OrderRepository.InsertOrderTx(ctx, tx, order)
		if err != nil {
			return err
		}
		orderID = order.ID
		return nil
	})
	if err != nil {
		return 0, err
	}
	return orderID, nil

}

func (s *OrderService) GetOrderHistoryByUserID(ctx context.Context, param *models.OrderHistoryParam) ([]models.OrderHistoryResponse, error) {
	orderHistories, err := s.OrderRepository.GetOrderHistoryByUserID(ctx, *param)
	if err != nil {
		return nil, err
	}
	return orderHistories, nil
}

func (s *OrderService) GetProductInfo(ctx context.Context, productID int64) (models.Product, error) {
	productInfo, err := s.OrderRepository.GetProductInfo(ctx, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Product{}, nil
		}
		return models.Product{}, err
	}
	return productInfo, nil
}

package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"order/order/cmd/order/service"
	constant "order/order/infrastructure/constans"
	"order/order/infrastructure/log"
	"order/order/kafka"
	"order/order/models"
	"time"

	"github.com/sirupsen/logrus"
)

type OrderUseCase struct {
	OrderService  *service.OrderService
	KafkaProducer kafka.KafkaProducer
}

func NewOrderUseCase(orderService *service.OrderService, kafkaProducer kafka.KafkaProducer) *OrderUseCase {
	return &OrderUseCase{
		OrderService:  orderService,
		KafkaProducer: kafkaProducer,
	}
}

func (uc *OrderUseCase) CheckOutOrder(ctx context.Context, param *models.CheckOutRequest) (int64, error) {
	var orderID int64

	//check idempotency
	if param.IdempotencyToken != "" {
		isExists, err := uc.OrderService.CheckIdempotency(ctx, param.IdempotencyToken)
		if err != nil {
			return 0, err
		}
		if isExists {
			return 0, errors.New("order already created, please check again")
		}
	}
	//validate product
	err := uc.validateProducsts(ctx, param.Items)
	if err != nil {
		return 0, err
	}
	//product amount
	totalQty, totalAmount := uc.calculateOrderSummary(param.Items)

	//construct order detail
	products, orderHistory := uc.constructOrderDetail(param.Items)

	//save order and detail
	orderDetail := models.OrderDetail{
		Products:     products,
		OrderHistory: orderHistory,
	}

	order := models.Order{
		UserID:          param.UserID,
		Amount:          totalAmount,
		TotalQty:        totalQty,
		Status:          constant.OrderStatusCreated,
		PaymentMethod:   param.PaymentMethod,
		ShippingAddress: param.ShippingAddress,
	}
	orderID, err = uc.OrderService.SaveOrderAndOrderDetail(ctx, &order, &orderDetail)
	if err != nil {
		return 0, err
	}

	if param.IdempotencyToken != "" {
		err = uc.OrderService.SaveIdempotency(ctx, param.IdempotencyToken)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"err":   err.Error(),
				"token": param.IdempotencyToken,
			}).Info("uc.OrderService.SaveIdempotency got error")
			return 0, err
		}
	}
	//publish order created
	orderCreatedEvent := models.OrderCreatedEvent{
		OrderID:         orderID,
		UserID:          param.UserID,
		PaymentMethod:   param.PaymentMethod,
		TotalAmount:     order.Amount,
		ShippingAddress: param.ShippingAddress,
	}

	err = uc.KafkaProducer.PublishOrderCreated(ctx, orderCreatedEvent)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

func (uc *OrderUseCase) validateProducsts(ctx context.Context, items []models.CheckOutItem) error {
	seen := map[int64]bool{}
	for _, item := range items {
		productInfo, err := uc.OrderService.GetProductInfo(ctx, item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to get product info for product %d: %v", item.ProductID, err)
		}

		if productInfo.ID == 0 {
			return fmt.Errorf("product with id %d not found", item.ProductID)
		}

		if item.Price != productInfo.Price {
			return fmt.Errorf("price mismatch for product %d: expected %f, got %f", item.ProductID, productInfo.Price, item.Price)
		}

		if seen[item.ProductID] {
			return fmt.Errorf("duplicate product: %d", item.ProductID)
		}
		seen[item.ProductID] = true
		if item.Quantity <= 0 || item.Quantity > 1000 {
			return fmt.Errorf("invalid quantity for product %d, maximum product qty is 100", item.ProductID)
		}

		if item.Price <= 0 {
			return fmt.Errorf("invalid price for product %d", item.ProductID)
		}

		if item.Quantity > productInfo.Stock {
			return fmt.Errorf("insufficient stock for product %d: available %d, requested %d", item.ProductID, productInfo.Stock, item.Quantity)
		}
	}
	return nil
}

func (uc *OrderUseCase) calculateOrderSummary(items []models.CheckOutItem) (int, float64) {
	var totalQty int
	var totalAmount float64

	for _, item := range items {
		totalQty += item.Quantity
		totalAmount += float64(item.Quantity) * item.Price
	}
	return totalQty, totalAmount
}

func (uc *OrderUseCase) constructOrderDetail(items []models.CheckOutItem) (string, string) {
	productJson, _ := json.Marshal(items)

	history := []map[string]interface{}{
		{"status": "created", "timestamp": time.Now()},
	}
	historyJSON, _ := json.Marshal(history)
	return string(productJson), string(historyJSON)

}

func (uc *OrderUseCase) GetOrderHistoryByUserID(ctx context.Context, param *models.OrderHistoryParam) ([]models.OrderHistoryResponse, error) {
	orderHistories, err := uc.OrderService.GetOrderHistoryByUserID(ctx, param)
	if err != nil {
		return nil, err
	}
	return orderHistories, nil
}

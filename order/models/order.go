package models

import "time"

type Order struct {
	ID              int64
	UserID          int64
	OrderDetailID   int64
	Amount          float64
	TotalQty        int
	Status          int
	PaymentMethod   string
	ShippingAddress string
}

type OrderDetail struct {
	ID           int64
	Products     string
	OrderHistory string
}

type CheckOutItem struct {
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
type CheckOutRequest struct {
	UserID           int64          `json:"user_id"`
	Items            []CheckOutItem `json:"items"`
	PaymentMethod    string         `json:"payment_method"`
	ShippingAddress  string         `json:"shipping_address"`
	IdompotencyToken string         `json:"idompotency_token"`
}

type OrderHistoryParam struct {
	UserID int64
	Status int
}

type OrderHistoryResponse struct {
	OrderID         int64           `json:"order_id"`
	TotalAmount     float64         `json:"total_amount"`
	TotalQty        int             `json:"total_qty"`
	Status          string          `json:"status"`
	PaymentMethod   string          `json:"payment_method"`
	ShippingAddress string          `json:"shipping_address"`
	Products        []CheckOutItem  `json:"products"`
	History         []Statushistory `json:"history"`
}

type Statushistory struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

type OrderRequestLog struct {
	ID               int64     `json:"id"`
	IdempotencyToken string    `json:"idempotancy_token"`
	CreateTime       time.Time `json:"create_time"`
}

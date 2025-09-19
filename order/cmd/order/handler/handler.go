package handler

import (
	"net/http"
	"order/order/cmd/order/usecase"
	"order/order/infrastructure/log"
	"order/order/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type OrderHandler struct {
	OrderUsecase *usecase.OrderUseCase
}

func NewOrderHandler(orderUseCase *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		OrderUsecase: orderUseCase,
	}
}

func (h *OrderHandler) CheckoutOrder(c *gin.Context) {
	var param models.CheckOutRequest
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "error_detail": err.Error()})
		return
	}
	userIDStr, isExist := c.Get("user_id")
	if !isExist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, ok := userIDStr.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id"})
		return
	}
	if len(param.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameter"})
		return
	}

	param.UserID = int64(userID)
	orderID, err := h.OrderUsecase.CheckOutOrder(c.Request.Context(), &param)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"param": param,
		}).Errorf("h.OrderUsecase.CheckOutOrder got error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Succesfull Create A Order",
		"orderID": orderID,
	})

}

func (h *OrderHandler) GetOrderHistory(c *gin.Context) {
	userIDStr, isExist := c.Get("user_id")
	if !isExist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, ok := userIDStr.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id"})
		return
	}
	statusStr := c.DefaultQuery("status", "0")
	status, _ := strconv.Atoi(statusStr)

	param := models.OrderHistoryParam{
		UserID: int64(userID),
		Status: status,
	}
	orderHistory, err := h.OrderUsecase.GetOrderHistoryByUserID(c.Request.Context(), &param)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"param": param,
		}).Errorf("h.OrderUsecase.GetOrderHistoryByUserID got error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": orderHistory,
	})

}

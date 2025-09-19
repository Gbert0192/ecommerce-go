package routes

import (
	"order/order/cmd/order/handler"
	"order/order/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, orderHander handler.OrderHandler, jwtSecret string) {
	router.Use(middleware.RequestLogger())
	router.Use(middleware.AuthMiddleware(jwtSecret))
	router.POST("/v1/checkout", orderHander.CheckoutOrder)
	router.GET("/v1/order_history", orderHander.GetOrderHistory)
}

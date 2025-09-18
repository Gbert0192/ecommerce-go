package routes

import (
	"order/order/cmd/order/handler"
	"order/order/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, orderHander handler.OrderHandler, jwtSecret string) {
	router.Use(middleware.RequestLogger())
	router.Use(middleware.AuthMiddleware(jwtSecret))
}

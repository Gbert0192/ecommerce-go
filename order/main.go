package main

import (
	"order/order/cmd/order/handler"
	"order/order/cmd/order/repository"
	"order/order/cmd/order/resource"
	"order/order/cmd/order/service"
	"order/order/cmd/order/usecase"
	"order/order/config"
	"order/order/infrastructure/log"
	routes "order/order/router"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	db := resource.InitDB(&cfg)
	redis := resource.InitRedis(&cfg)

	port := cfg.App.Port
	router := gin.Default()
	orderRepository := repository.NewOrderRepository(db, redis)
	orderService := service.NewOrderService(orderRepository)
	orderUseCase := usecase.NewOrderUseCase(orderService)
	orderHandler := handler.NewOrderHandler(orderUseCase)
	routes.SetupRouter(router, *orderHandler, cfg.Secret.JWTSecret)

	println("Starting server on port " + port)

	log.SetupLoger()
	router.Run(":" + port)
	println("Starting server on port:", port)

	log.Logger.Printf("Server Running on port: %s", port)

}

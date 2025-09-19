package main

import (
	"order/order/cmd/order/handler"
	"order/order/cmd/order/repository"
	"order/order/cmd/order/resource"
	"order/order/cmd/order/service"
	"order/order/cmd/order/usecase"
	"order/order/config"
	"order/order/infrastructure/log"
	"order/order/kafka"
	routes "order/order/router"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	db := resource.InitDB(&cfg)
	redis := resource.InitRedis(&cfg)
	KafkaProducer := kafka.NewKafkaProducer([]string{"localhost:9093"}, "order.created")
	defer KafkaProducer.Close()

	port := cfg.App.Port
	router := gin.Default()
	orderRepository := repository.NewOrderRepository(db, redis, cfg.Product.Host)
	orderService := service.NewOrderService(orderRepository)
	orderUseCase := usecase.NewOrderUseCase(orderService, *KafkaProducer)
	orderHandler := handler.NewOrderHandler(orderUseCase)
	routes.SetupRouter(router, *orderHandler, cfg.Secret.JWTSecret)

	println("Starting server on port " + port)

	log.SetupLoger()
	router.Run(":" + port)
	println("Starting server on port:", port)

	log.Logger.Printf("Server Running on port: %s", port)

}

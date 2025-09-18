package main

import (
	"product/cmd/product/handler"
	"product/cmd/product/repository"
	"product/cmd/product/resource"
	"product/cmd/product/service"
	"product/cmd/product/usecase"
	"product/config"
	"product/infrastructure/log"
	routes "product/router"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	redis := resource.InitRedis(&cfg)
	db := resource.InitDB(&cfg)

	log.SetupLoger()

	productRepository := repository.NewProductRepository(db, redis)
	productService := service.NewProductRepository(*productRepository)
	productUseCase := usecase.NewProductUseCase(*productService)
	productHandler := handler.NewProductHandler(productUseCase)
	port := cfg.App.Port
	router := gin.Default()

	routes.SetupRouter(router, *productHandler)

	router.Run(":" + port)
	println("Starting server on port:", port)

	log.Logger.Printf("Server Running on port: %s", port)

}

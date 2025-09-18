package main

import (
	"user/cmd/user/handler"
	"user/cmd/user/repository"
	"user/cmd/user/resource"
	"user/cmd/user/service"
	"user/cmd/user/usecase"
	"user/config"
	"user/infrastructure/log"
	"user/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	redis := resource.InitRedis(&cfg)
	db := resource.IntDb(&cfg)
	log.SetupLoger()

	userRepository := repository.NewUserRepository(redis, db)
	userService := service.NewUserService(*userRepository)
	userUseCase := usecase.NewUserUseCase(*userService, cfg.Secret.JWTSecret)
	UserHandler := handler.NewUserHandler(*userUseCase)

	port := cfg.App.Port
	router := gin.Default()

	routes.SetupRoutes(router, *UserHandler, cfg.Secret.JWTSecret)

	router.Run(":" + port)

	log.Logger.Printf("Server running on port %s", port)

}

package resource

import (
	"context"
	"fmt"
	"log"
	"user/config"

	"github.com/redis/go-redis/v9"
)

var RedistClient *redis.Client

func InitRedis(cfg *config.Config) *redis.Client {
	RedistClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port), //host:port
		Password: cfg.Redis.Password,
	})

	ctx := context.Background()
	pingResult, err := RedistClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	log.Println("Connected to redis:", pingResult)
	return RedistClient

}

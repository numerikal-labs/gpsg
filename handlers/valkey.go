package handlers

import (
	"log"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var (
	valkeyClient *redis.Client
	ctx          = context.Background()
)

func InitValkey() error {
	valkeyClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := valkeyClient.Ping(ctx).Result()
	if err != nil {
		return err
	}
	log.Println("Connected to Valkey")
	return nil
}

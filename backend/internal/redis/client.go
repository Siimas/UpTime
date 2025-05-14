package redis

import (
	"os"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func NewClient() *redis.Client {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}
	
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	return redis.NewClient(opt)
}
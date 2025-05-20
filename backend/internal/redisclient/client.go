package redisclient

import (
	"log"
	"github.com/redis/go-redis/v9"
)

func NewClient(url string) *redis.Client {	
	opt, err := redis.ParseURL(url)
	if err != nil {
		log.Fatalf("Couldn't establish connection with redis: %s\n", err)
	}
	return redis.NewClient(opt)
}
package redis

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func NewClient(url string) (*redis.Client, error) {	
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		return nil, err
	}
	return redis.NewClient(opt), nil
}
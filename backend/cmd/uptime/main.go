package main

import (
	"context"
	"uptime/internal/monitor"
	"uptime/internal/redis"
)

func main() {
	ctx := context.Background()
	rdb := redis.NewClient()
	monitor.RunScheduler(ctx, rdb)
}

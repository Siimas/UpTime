package monitor

import (
	"context"
	"fmt"
	"sync"
	"uptime/internal/util"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func RunMonitorFlusher(ctx context.Context, db *pgx.Conn, kc *kafka.Consumer, rdb *redis.Client) {
	if err := LoadMonitorConfigs(ctx, db, rdb); err != nil {
		fmt.Println("Error starting flusher: " + err.Error())
	}

	// for {

	// }
}

func LoadMonitorConfigs(ctx context.Context, db *pgx.Conn, rdb *redis.Client) error {
	monitors, err := GetActiveMonitors(ctx, db)
	if err != nil {
		fmt.Println("Failed to load monitors: " + err.Error())
	}

	var wg sync.WaitGroup

	for _, m := range monitors {
		wg.Add(1)

		go func(m Monitor) {
			defer wg.Done()
			if err := ScheduleMonitor(ctx, m, rdb); err != nil {
				fmt.Printf("Failed schedule monitor (%s): %s", m.Id, err.Error())
			}
			fmt.Println("Monitor Scheduled!")
			util.PrettyPrint(m)
		}(m)
	}

	wg.Wait()
	return nil
}

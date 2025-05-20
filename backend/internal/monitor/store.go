package monitor

import (
	"context"
	"uptime/internal/models"

	"github.com/jackc/pgx/v5"
)

func StoreMonitorResult(ctx context.Context, mr models.MonitorResult, db *pgx.Conn) error {
	sql := `INSERT INTO monitor_results (monitor_id, status, latency_ms, response_code, error, checked_at) VALUES ($1, $2, $3, $4, $5, $6)`

	if _, err := db.Exec(ctx, sql, mr.Id, mr.Status, mr.Latency, mr.Code, mr.Error, mr.Date); err != nil {
		return err
	}

	return nil
}


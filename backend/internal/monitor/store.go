package monitor

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func StoreMonitorResult(ctx context.Context, mr MonitorResult, db *pgx.Conn) error {
	sql := `INSERT INTO monitor_results (monitor_id, status, latency_ms, response_code, error, checked_at) VALUES ($1, $2, $3, $4, $5, $6)`

	if _, err := db.Exec(ctx, sql, mr.Id, mr.Status, mr.Latency, mr.Code, mr.Error, mr.Date); err != nil {
		return err
	}

	return nil
}

func GetActiveMonitors(ctx context.Context, db *pgx.Conn) ([]Monitor, error) {
	rows, err := db.Query(ctx, `
        SELECT id, url, interval_seconds, active
        FROM monitors
        WHERE active = true
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monitors []Monitor

	for rows.Next() {
		var m Monitor
		err := rows.Scan(&m.Id, &m.Endpoint, &m.Interval, &m.Active)
		if err != nil {
			return nil, err
		}
		monitors = append(monitors, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return monitors, nil
}

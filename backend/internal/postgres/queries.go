package postgres

import (
	"context"
	"uptime/internal/models"

	"github.com/jackc/pgx/v5"
)

func GetSingleMonitor(ctx context.Context, db *pgx.Conn, monitorId string) (models.Monitor, error) {
	row := db.QueryRow(ctx, "SELECT id, url, interval_seconds, active FROM monitors WHERE id = $1", monitorId)

	var m models.Monitor
	if err := row.Scan(&m.Id, &m.Endpoint, &m.Interval, &m.Active); err != nil {
		return models.Monitor{}, err
	}

	return m, nil
}

func GetActiveMonitors(ctx context.Context, db *pgx.Conn) ([]models.Monitor, error) {
	rows, err := db.Query(ctx, `
        SELECT id, url, interval_seconds, active
        FROM monitors
        WHERE active = true
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monitors []models.Monitor

	for rows.Next() {
		var m models.Monitor
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

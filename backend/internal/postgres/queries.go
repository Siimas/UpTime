package postgres

import (
	"context"
	"uptime/internal/models"

	"github.com/jackc/pgx/v5"
)

func GetAllMonitors(ctx context.Context, db *pgx.Conn) ([]models.Monitor, error) {
	rows, err := db.Query(ctx, "SELECT id, url, interval_seconds, active FROM monitors")
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

	return monitors, nil
}

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

func CreateMonitor(ctx context.Context, db *pgx.Conn, m models.MonitorCreateDTO) (string, error) {
	var monitorId string
	err := db.QueryRow(
		ctx,
		"INSERT INTO monitors (url, interval_seconds, active) VALUES ($1, $2, $3) RETURNING id",
		m.Endpoint, m.Interval, m.Active,
	).Scan(&monitorId)

	if err != nil {
		return "", err
	}

	return monitorId, nil
}

func UpdateMonitor(ctx context.Context, db *pgx.Conn, m models.MonitorUpdateDTO) (int64, error) {
	result, err := db.Exec(
		ctx,
		"UPDATE monitors SET url = $1, interval_seconds = $2, active = $3 WHERE id = $4",
		m.Endpoint, m.Interval, m.Active, m.Id,
	)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func DeleteMonitor(ctx context.Context, db *pgx.Conn, monitorId string) (int64, error) {
	result, err := db.Exec(
		ctx,
		"DELETE FROM monitors WHERE id = $1",
		monitorId,
	)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

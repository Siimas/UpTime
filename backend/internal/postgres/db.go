package postgres

import (
	"context"
	"log"
	"uptime/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConnection(ctx context.Context) *pgx.Conn {
	con, err := pgx.Connect(context.Background(), config.GetEnv("POSTGRES_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	return con
}

func NewPoolConnection(ctx context.Context, url string) *pgxpool.Pool {
	pcon, err := pgxpool.New(ctx, url)
	if err != nil {
		log.Fatalf("Unable to create database pool: %v\n", err)
	}
	return pcon
}
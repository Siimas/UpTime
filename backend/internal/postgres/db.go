package postgres

import (
	"log"
	"context"
	"github.com/jackc/pgx/v5"
)

func NewConnection(ctx context.Context, url string) *pgx.Conn {
	con, err := pgx.Connect(context.Background(), url)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	return con
}
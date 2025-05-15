package postgres

import (
	"os"
	"fmt"
	"context"

	"github.com/jackc/pgx/v5"
)

func NewConnection(ctx context.Context, url string) *pgx.Conn {
	con, err := pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return con
}
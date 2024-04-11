package model

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func Setup(ctx context.Context, DSN string) *pgx.Conn {
	conn, err := pgx.Connect(ctx, DSN)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

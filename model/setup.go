package model

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func Setup(ctx context.Context, DSN string) *Queries {
	conn, err := pgx.Connect(ctx, DSN)
	if err != nil {
		log.Fatal(err)
	}
	return New(conn)
}

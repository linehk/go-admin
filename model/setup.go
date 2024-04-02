package model

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/linehk/go-admin/config"
)

var Query *Queries

func Setup() {
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.Raw.String("HOST"), config.Raw.String("POSTGRES_USER"), config.Raw.String("POSTGRES_PASSWORD"),
		config.Raw.String("POSTGRES_DB"), config.Raw.String("PORT"), config.Raw.String("SSL_MODE"),
		config.Raw.String("TIMEZONE"))

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, DSN)
	if err != nil {
		log.Fatalf("can't open database err: %v", err)
	}
	Query = New(conn)
}

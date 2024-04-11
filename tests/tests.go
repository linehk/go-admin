package tests

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/linehk/go-admin/model"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const BaseURL = "http://localhost:8080/"

func ContainerDB(t *testing.T) *pgx.Conn {
	ctx := context.Background()
	pg, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16.2"),
		postgres.WithInitScripts(filepath.Join("..", "..", "model", "schema.sql")),
		postgres.WithDatabase("go_admin"),
		postgres.WithUsername("dev"),
		postgres.WithPassword("dev"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)))
	if err != nil {
		t.Error(err)
	}
	dsn, err := pg.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Error(err)
	}
	return model.Setup(ctx, dsn)
}

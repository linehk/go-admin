package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/linehk/go-admin/controller"
	"github.com/linehk/go-admin/model"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestPostApiV1Users(t *testing.T) {
	reqBody := "{\"password\":\"password\",\"username\":\"username\"}\n"
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/users", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	ctx := context.Background()
	pg, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16.2"),
		postgres.WithInitScripts(filepath.Join("..", "model", "schema.sql")),
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

	userImpl := &controller.UserImpl{DB: model.Setup(ctx, dsn)}

	userImpl.PostApiV1Users(w, req)
	actual := w.Body.String()
	expected := "{\"password\":\"password\",\"username\":\"username\"}\n"
	assert.Equal(t, expected, actual)
}

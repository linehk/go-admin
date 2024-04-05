package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/linehk/go-admin/config"
	"github.com/linehk/go-admin/model"
	"golang.org/x/crypto/bcrypt"
)

func Setup() *http.ServeMux {
	mux := http.NewServeMux()
	ctx := context.Background()
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.Raw.String("HOST"), config.Raw.String("POSTGRES_USER"), config.Raw.String("POSTGRES_PASSWORD"),
		config.Raw.String("POSTGRES_DB"), config.Raw.String("PORT"), config.Raw.String("SSL_MODE"),
		config.Raw.String("TIMEZONE"))
	userImpl := &UserImpl{DB: model.Setup(ctx, DSN)}
	options := StdHTTPServerOptions{
		BaseRouter: mux,
	}
	HandlerWithOptions(userImpl, options)
	return mux
}

func hash(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(h), err
}

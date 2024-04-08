package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/linehk/go-admin/config"
	"github.com/linehk/go-admin/errcode"
	"github.com/linehk/go-admin/model"
)

type API struct {
	DB *model.Queries
}

func Setup() *http.ServeMux {
	mux := http.NewServeMux()
	ctx := context.Background()
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.Raw.String("HOST"), config.Raw.String("POSTGRES_USER"), config.Raw.String("POSTGRES_PASSWORD"),
		config.Raw.String("POSTGRES_DB"), config.Raw.String("PORT"), config.Raw.String("SSL_MODE"),
		config.Raw.String("TIMEZONE"))
	api := &API{DB: model.Setup(ctx, DSN)}
	options := StdHTTPServerOptions{
		BaseRouter: mux,
	}
	HandlerWithOptions(api, options)
	return mux
}

func decode(w http.ResponseWriter, r *http.Request, req any) {
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		returnErr(w, errcode.Parse)
	}
}

func encode(w http.ResponseWriter, resp any) {
	err := json.NewEncoder(w).Encode(&resp)
	if err != nil {
		returnErr(w, errcode.Parse)
	}
}

func paging(current, pageSize int32) (int32, int32) {
	if current > 0 && pageSize > 0 {
		return (current - 1) * pageSize, pageSize
	}
	return current, pageSize
}

func returnErr(w http.ResponseWriter, e int32) {
	errResp := Error{
		Code:    e,
		Message: errcode.Msg(e),
	}
	err := json.NewEncoder(w).Encode(errResp)
	if err != nil {
		panic(err)
	}
	slog.Error("err: ", err)
	w.WriteHeader(http.StatusBadRequest)
}

const pgTimestampFormat = "2006-01-02 15:04:05.999999999"

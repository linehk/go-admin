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
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ReturnErr(w, errcode.Parse)
	}
}

func encode(w http.ResponseWriter, resp any) {
	err := json.NewEncoder(w).Encode(&resp)
	if err != nil {
		ReturnErr(w, errcode.Parse)
	}
}

func paging(current, pageSize int) (int32, int32) {
	if current > 0 && pageSize > 0 {
		return int32((current - 1) * pageSize), int32(pageSize)
	}
	return int32(current), int32(pageSize)
}

func ReturnErr(w http.ResponseWriter, e int32) {
	w.Header().Set("Content-Type", "application/json")
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

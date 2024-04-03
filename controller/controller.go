package controller

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/linehk/go-admin/model"
)

func Setup() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("POST /api/users", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user model.AppUser
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			slog.Error("err: ", err)
			return
		}

		ctx := context.Background()

		userResult, err := model.Query.CreateUser(ctx, model.CreateUserParams{
			Username: user.Username,
			Password: user.Password,
		})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			slog.Error("err: ", err)
			return
		}

		err = json.NewEncoder(w).Encode(userResult)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			slog.Error("err: ", err)
			return
		}
	}))

	return mux
}

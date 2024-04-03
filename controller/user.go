package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/linehk/go-admin/model"
)

type UserImpl struct {
	DB *model.Queries
}

func (u *UserImpl) PostApiV1Users(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user PostApiV1UsersJSONRequestBody
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "decode err: ", http.StatusBadRequest)
		slog.Error("decode err: ", err)
		return
	}

	userModel, err := u.DB.CreateUser(r.Context(), model.CreateUserParams{
		Username: *user.Username,
		Password: *user.Password,
	})
	if err != nil {
		http.Error(w, "db err: ", http.StatusBadRequest)
		slog.Error("db err: ", err)
		return
	}

	var userResp User
	userResp.Username = &userModel.Username
	userResp.Password = &userModel.Password

	err = json.NewEncoder(w).Encode(userResp)
	if err != nil {
		http.Error(w, "decode err: ", http.StatusBadRequest)
		slog.Error("decode err: ", err)
		return
	}
}

func (u *UserImpl) GetApiV1Users(w http.ResponseWriter, r *http.Request, params GetApiV1UsersParams) {
}

func (u *UserImpl) DeleteApiV1UsersId(w http.ResponseWriter, r *http.Request, id int64) {}

func (u *UserImpl) GetApiV1UsersId(w http.ResponseWriter, r *http.Request, id int64) {}

package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/linehk/go-admin/model"
)

type UserImpl struct {
}

func (*UserImpl) PostApiV1Users(w http.ResponseWriter, r *http.Request) {
	var user PostApiV1UsersJSONRequestBody
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "decode err: ", http.StatusBadRequest)
		slog.Error("decode err: ", err)
		return
	}

	userModel, err := model.DB.CreateUser(r.Context(), model.CreateUserParams{
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

func (*UserImpl) GetApiV1Users(w http.ResponseWriter, r *http.Request, params GetApiV1UsersParams) {}

func (*UserImpl) DeleteApiV1UsersId(w http.ResponseWriter, r *http.Request, id int64) {}

func (*UserImpl) GetApiV1UsersId(w http.ResponseWriter, r *http.Request, id int64) {}

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
	var req PostApiV1UsersJSONRequestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "decode err: ", http.StatusBadRequest)
		slog.Error("decode err: ", err)
		return
	}

	createUserParams, err := reqToCreateUserParams(req)
	if err != nil {
		http.Error(w, "reqToCreateUserParams err: ", http.StatusBadRequest)
		slog.Error("reqToCreateUserParams err: ", err)
		return
	}

	userModel, err := u.DB.CreateUser(r.Context(), createUserParams)
	if err != nil {
		http.Error(w, "db err: ", http.StatusBadRequest)
		slog.Error("db err: ", err)
		return
	}

	err = json.NewEncoder(w).Encode(userModelToResp(userModel))
	if err != nil {
		http.Error(w, "decode err: ", http.StatusBadRequest)
		slog.Error("decode err: ", err)
		return
	}
}

func (u *UserImpl) GetApiV1Users(w http.ResponseWriter, r *http.Request, params GetApiV1UsersParams) {
}

func (u *UserImpl) DeleteApiV1UsersId(w http.ResponseWriter, r *http.Request, id int32) {
	w.Header().Set("Content-Type", "application/json")
	err := u.DB.DeleteUser(r.Context(), id)
	if err != nil {
		http.Error(w, "db err: ", http.StatusBadRequest)
		slog.Error("db err: ", err)
		return
	}
}

func (u *UserImpl) GetApiV1UsersId(w http.ResponseWriter, r *http.Request, id int32) {
	w.Header().Set("Content-Type", "application/json")
	userModel, err := u.DB.GetUser(r.Context(), id)
	if err != nil {
		http.Error(w, "db err: ", http.StatusBadRequest)
		slog.Error("db err: ", err)
		return
	}
	err = json.NewEncoder(w).Encode(userModelToResp(userModel))
	if err != nil {
		http.Error(w, "decode err: ", http.StatusBadRequest)
		slog.Error("decode err: ", err)
		return
	}
}

func (u *UserImpl) PutApiV1UsersId(w http.ResponseWriter, r *http.Request, id int32) {
	w.Header().Set("Content-Type", "application/json")
	var req PutApiV1UsersIdJSONRequestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "decode err: ", http.StatusBadRequest)
		slog.Error("decode err: ", err)
		return
	}
	updateUserParams, err := reqToUpdateUserParams(req)
	if err != nil {
		http.Error(w, "reqToUpdateUserParams err: ", http.StatusBadRequest)
		slog.Error("reqToUpdateUserParams err: ", err)
		return
	}
	updateUserParams.ID = id
	userModel, err := u.DB.UpdateUser(r.Context(), updateUserParams)
	if err != nil {
		http.Error(w, "db err: ", http.StatusBadRequest)
		slog.Error("db err: ", err)
		return
	}

	err = json.NewEncoder(w).Encode(userModelToResp(userModel))
	if err != nil {
		http.Error(w, "decode err: ", http.StatusBadRequest)
		slog.Error("decode err: ", err)
		return
	}
}

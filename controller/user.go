package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

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

	var createUserParams model.CreateUserParams
	createUserParams.Username = user.Username
	password, err := hash(user.Password)
	if err != nil {
		http.Error(w, "hash password err: ", http.StatusBadRequest)
		slog.Error("hash password err: ", err)
		return
	}
	createUserParams.Password = password
	if user.Name != nil {
		createUserParams.Name = *user.Name
	}
	if user.Email != nil {
		createUserParams.Email = *user.Email
	}
	if user.Phone != nil {
		createUserParams.Phone = *user.Phone
	}
	if user.Remark != nil {
		createUserParams.Remark = *user.Remark
	}
	createUserParams.Status = string(Activated)
	if user.Created != nil {
		err = createUserParams.Created.Scan(*user.Created)
		if err != nil {
			http.Error(w, "scan time err: ", http.StatusBadRequest)
			slog.Error("scan time err: ", err)
			return
		}
	} else {
		err = createUserParams.Created.Scan(time.Now())
		if err != nil {
			http.Error(w, "scan time err: ", http.StatusBadRequest)
			slog.Error("scan time err: ", err)
			return
		}
	}
	if user.Updated != nil {
		err = createUserParams.Updated.Scan(*user.Updated)
		if err != nil {
			http.Error(w, "scan time err: ", http.StatusBadRequest)
			slog.Error("scan time err: ", err)
			return
		}
	} else {
		err = createUserParams.Updated.Scan(time.Now())
		if err != nil {
			http.Error(w, "scan time err: ", http.StatusBadRequest)
			slog.Error("scan time err: ", err)
			return
		}
	}

	userModel, err := u.DB.CreateUser(r.Context(), createUserParams)
	if err != nil {
		http.Error(w, "db err: ", http.StatusBadRequest)
		slog.Error("db err: ", err)
		return
	}

	err = json.NewEncoder(w).Encode(CreateUserRowToResp(userModel))
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
	err = json.NewEncoder(w).Encode(GetUserRowToResp(userModel))
	if err != nil {
		http.Error(w, "decode err: ", http.StatusBadRequest)
		slog.Error("decode err: ", err)
		return
	}
}

func (u *UserImpl) PutApiV1UsersId(w http.ResponseWriter, r *http.Request, id int32) {}

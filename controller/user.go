package controller

import (
	"net/http"

	"github.com/linehk/go-admin/errcode"
	"github.com/linehk/go-admin/model"
)

type UserImpl struct {
	DB *model.Queries
}

func (u *UserImpl) PostApiV1Users(w http.ResponseWriter, r *http.Request) {
	var req User
	decode(w, r, &req)

	createUserParams, err := reqToCreateUserParams(req)
	if err != nil {
		ReturnErr(w, errcode.Convert)
	}

	userModel, err := u.DB.CreateUser(r.Context(), createUserParams)
	if err != nil {
		ReturnErr(w, errcode.Database)
	}

	userResp := userModelToResp(userModel)
	encode(w, userResp)
}

func (u *UserImpl) GetApiV1Users(w http.ResponseWriter, r *http.Request, params GetApiV1UsersParams) {
	w.Header().Set("Content-Type", "application/json")
	var listUserParams model.ListUserParams
	if params.Username != nil {
		listUserParams.Column1 = *params.Username
	}
	if params.Name != nil {
		listUserParams.Column2 = *params.Name
	}
	if params.Status != nil {
		listUserParams.Column3 = *params.Status
	}

	listUserParams.ID, listUserParams.Limit = paging(params.Current, params.PageSize)

	userModelList, err := u.DB.ListUser(r.Context(), listUserParams)
	if err != nil {
		ReturnErr(w, errcode.Database)
	}

	var userRespList []User
	for _, userModel := range userModelList {
		userRespList = append(userRespList, userModelToResp(userModel))
	}

	encode(w, userRespList)
}

func (u *UserImpl) DeleteApiV1UsersId(w http.ResponseWriter, r *http.Request, id int32) {
	w.Header().Set("Content-Type", "application/json")
	err := u.DB.DeleteUser(r.Context(), id)
	if err != nil {
		ReturnErr(w, errcode.Database)
	}
}

func (u *UserImpl) GetApiV1UsersId(w http.ResponseWriter, r *http.Request, id int32) {
	w.Header().Set("Content-Type", "application/json")
	userModel, err := u.DB.GetUser(r.Context(), id)
	if err != nil {
		ReturnErr(w, errcode.Database)
	}

	userResp := userModelToResp(userModel)
	encode(w, userResp)
}

func (u *UserImpl) PutApiV1UsersId(w http.ResponseWriter, r *http.Request, id int32) {
	var req User
	decode(w, r, &req)

	updateUserParams, err := reqToUpdateUserParams(req)
	if err != nil {
		ReturnErr(w, errcode.Convert)
	}

	updateUserParams.ID = id
	userModel, err := u.DB.UpdateUser(r.Context(), updateUserParams)
	if err != nil {
		ReturnErr(w, errcode.Database)
	}

	userResp := userModelToResp(userModel)
	encode(w, userResp)
}

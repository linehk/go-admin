package controller

import (
	"net/http"
	"time"

	"github.com/linehk/go-admin/errcode"
	"github.com/linehk/go-admin/model"
	"golang.org/x/crypto/bcrypt"
)

func (a *API) PostApiV1Users(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req User
	decode(w, r, &req)

	params, err := createUserParams(req)
	if err != nil {
		Err(w, errcode.Convert)
		return
	}

	user, err := a.DB.CreateUser(r.Context(), params)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	resp := userResp(user)
	encode(w, resp)
}

func (a *API) GetApiV1Users(w http.ResponseWriter, r *http.Request, params GetApiV1UsersParams) {
	w.Header().Set("Content-Type", "application/json")

	var modelParams model.ListUserParams
	modelParams.Column1 = params.Username
	modelParams.Column2 = params.Name
	modelParams.Column3 = params.Status
	modelParams.ID, modelParams.Limit = paging(params.Current, params.PageSize)

	userList, err := a.DB.ListUser(r.Context(), modelParams)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	var respList []User
	for _, userModel := range userList {
		respList = append(respList, userResp(userModel))
	}

	encode(w, respList)
}

func (a *API) DeleteApiV1UsersId(w http.ResponseWriter, r *http.Request, id int32) {
	w.Header().Set("Content-Type", "application/json")

	err := a.DB.DeleteUser(r.Context(), id)
	if err != nil {
		Err(w, errcode.Database)
		return
	}
}

func (a *API) GetApiV1UsersId(w http.ResponseWriter, r *http.Request, id int32) {
	w.Header().Set("Content-Type", "application/json")

	user, err := a.DB.GetUser(r.Context(), id)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	resp := userResp(user)
	encode(w, resp)
}

func (a *API) PutApiV1UsersId(w http.ResponseWriter, r *http.Request, id int32) {
	w.Header().Set("Content-Type", "application/json")

	var req User
	decode(w, r, &req)

	params, err := updateUserParams(req)
	if err != nil {
		Err(w, errcode.Convert)
		return
	}

	params.ID = id
	user, err := a.DB.UpdateUser(r.Context(), params)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	resp := userResp(user)
	encode(w, resp)
}

func createUserParams(req User) (model.CreateUserParams, error) {
	var params model.CreateUserParams
	params.Username = req.Username
	password, err := hash(req.Password)
	if err != nil {
		return model.CreateUserParams{}, err
	}
	params.Password = password
	params.Name = req.Name
	params.Email = req.Email
	params.Phone = req.Phone
	params.Remark = req.Remark
	if req.Status == "" {
		req.Status = Activated
	}
	params.Status = string(req.Status)
	if req.Created == "" {
		req.Created = time.Now().Format(pgTimestampFormat)
	}
	err = params.Created.Scan(req.Created)
	if err != nil {
		return model.CreateUserParams{}, err
	}
	if req.Updated == "" {
		req.Updated = time.Now().Format(pgTimestampFormat)
	}
	err = params.Updated.Scan(req.Updated)
	if err != nil {
		return model.CreateUserParams{}, err
	}
	return params, nil
}

func updateUserParams(req User) (model.UpdateUserParams, error) {
	var params model.UpdateUserParams
	params.Username = req.Username
	password, err := hash(req.Password)
	if err != nil {
		return model.UpdateUserParams{}, err
	}
	params.Password = password
	params.Name = req.Name
	params.Email = req.Email
	params.Phone = req.Phone
	params.Remark = req.Remark
	params.Status = string(req.Status)
	err = params.Created.Scan(req.Created)
	if err != nil {
		return model.UpdateUserParams{}, err
	}
	if req.Updated == "" {
		req.Updated = time.Now().Format(pgTimestampFormat)
	}
	err = params.Updated.Scan(req.Updated)
	if err != nil {
		return model.UpdateUserParams{}, err
	}
	return params, nil
}

func hash(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(h), err
}

func userResp(userModel model.AppUser) User {
	var resp User
	resp.Id = &userModel.ID
	resp.Username = userModel.Username
	resp.Name = userModel.Name
	resp.Email = userModel.Email
	resp.Phone = userModel.Phone
	resp.Remark = userModel.Remark
	resp.Status = UserStatus(userModel.Status)
	resp.Created = userModel.Created.Time.Format(pgTimestampFormat)
	resp.Updated = userModel.Updated.Time.Format(pgTimestampFormat)
	return resp
}

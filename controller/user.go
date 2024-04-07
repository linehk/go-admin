package controller

import (
	"net/http"
	"time"

	"github.com/linehk/go-admin/errcode"
	"github.com/linehk/go-admin/model"
	"golang.org/x/crypto/bcrypt"
)

func (a *API) PostApiV1Users(w http.ResponseWriter, r *http.Request) {
	var req User
	decode(w, r, &req)

	createUserParams, err := reqToCreateUserParams(req)
	if err != nil {
		ReturnErr(w, errcode.Convert)
	}

	userModel, err := a.DB.CreateUser(r.Context(), createUserParams)
	if err != nil {
		ReturnErr(w, errcode.Database)
	}

	userResp := userModelToResp(userModel)
	encode(w, userResp)
}

func (a *API) GetApiV1Users(w http.ResponseWriter, r *http.Request, params GetApiV1UsersParams) {
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

	userModelList, err := a.DB.ListUser(r.Context(), listUserParams)
	if err != nil {
		ReturnErr(w, errcode.Database)
	}

	var userRespList []User
	for _, userModel := range userModelList {
		userRespList = append(userRespList, userModelToResp(userModel))
	}

	encode(w, userRespList)
}

func (a *API) DeleteApiV1UsersId(w http.ResponseWriter, r *http.Request, id int32) {
	w.Header().Set("Content-Type", "application/json")
	err := a.DB.DeleteUser(r.Context(), id)
	if err != nil {
		ReturnErr(w, errcode.Database)
	}
}

func (a *API) GetApiV1UsersId(w http.ResponseWriter, r *http.Request, id int32) {
	w.Header().Set("Content-Type", "application/json")
	userModel, err := a.DB.GetUser(r.Context(), id)
	if err != nil {
		ReturnErr(w, errcode.Database)
	}

	userResp := userModelToResp(userModel)
	encode(w, userResp)
}

func (a *API) PutApiV1UsersId(w http.ResponseWriter, r *http.Request, id int32) {
	var req User
	decode(w, r, &req)

	updateUserParams, err := reqToUpdateUserParams(req)
	if err != nil {
		ReturnErr(w, errcode.Convert)
	}

	updateUserParams.ID = id
	userModel, err := a.DB.UpdateUser(r.Context(), updateUserParams)
	if err != nil {
		ReturnErr(w, errcode.Database)
	}

	userResp := userModelToResp(userModel)
	encode(w, userResp)
}

func hash(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(h), err
}

func userModelToResp(userModel model.AppUser) User {
	const format = "2006-01-02 15:04:05.999999999"
	var userResp User
	userResp.Username = userModel.Username
	userResp.Name = &userModel.Name
	userResp.Email = &userModel.Email
	userResp.Phone = &userModel.Phone
	userResp.Remark = &userModel.Remark
	userStatus := UserStatus(userModel.Status)
	userResp.Status = &userStatus
	createdStr := userModel.Created.Time.Format(format)
	userResp.Created = &createdStr
	updatedStr := userModel.Updated.Time.Format(format)
	userResp.Updated = &updatedStr
	return userResp
}

func reqToCreateUserParams(req User) (model.CreateUserParams, error) {
	var createUserParams model.CreateUserParams
	createUserParams.Username = req.Username
	password, err := hash(req.Password)
	if err != nil {
		return model.CreateUserParams{}, err
	}
	createUserParams.Password = password
	if req.Name != nil {
		createUserParams.Name = *req.Name
	}
	if req.Email != nil {
		createUserParams.Email = *req.Email
	}
	if req.Phone != nil {
		createUserParams.Phone = *req.Phone
	}
	if req.Remark != nil {
		createUserParams.Remark = *req.Remark
	}
	if req.Status != nil {
		createUserParams.Status = string(*req.Status)
	} else {
		createUserParams.Status = string(Activated)
	}
	if req.Created != nil {
		err := createUserParams.Created.Scan(*req.Created)
		if err != nil {
			return model.CreateUserParams{}, err
		}
	} else {
		err := createUserParams.Created.Scan(time.Now())
		if err != nil {
			return model.CreateUserParams{}, err
		}
	}
	if req.Updated != nil {
		err := createUserParams.Updated.Scan(*req.Updated)
		if err != nil {
			return model.CreateUserParams{}, err
		}
	} else {
		err := createUserParams.Updated.Scan(time.Now())
		if err != nil {
			return model.CreateUserParams{}, err
		}
	}
	return createUserParams, nil
}

func reqToUpdateUserParams(req User) (model.UpdateUserParams, error) {
	var updateUserParams model.UpdateUserParams
	updateUserParams.Username = req.Username
	password, err := hash(req.Password)
	if err != nil {
		return model.UpdateUserParams{}, err
	}
	updateUserParams.Password = password
	if req.Name != nil {
		updateUserParams.Name = *req.Name
	}
	if req.Email != nil {
		updateUserParams.Email = *req.Email
	}
	if req.Phone != nil {
		updateUserParams.Phone = *req.Phone
	}
	if req.Remark != nil {
		updateUserParams.Remark = *req.Remark
	}
	if req.Status != nil {
		updateUserParams.Status = string(*req.Status)
	}
	if req.Created != nil {
		err := updateUserParams.Created.Scan(*req.Created)
		if err != nil {
			return model.UpdateUserParams{}, err
		}
	}
	if req.Updated != nil {
		err := updateUserParams.Updated.Scan(*req.Updated)
		if err != nil {
			return model.UpdateUserParams{}, err
		}
	} else {
		err := updateUserParams.Updated.Scan(time.Now())
		if err != nil {
			return model.UpdateUserParams{}, err
		}
	}
	return updateUserParams, nil
}

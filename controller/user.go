package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/linehk/go-admin/errcode"
	"github.com/linehk/go-admin/model"
	"golang.org/x/crypto/bcrypt"
)

func (a *API) PostApiV1Users(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	var req User
	decode(w, r, &req)
	err := validator.New().Struct(req)
	if err != nil {
		Err(w, errcode.Validate)
		return
	}

	params, err := createUserParams(req)
	if err != nil {
		Err(w, errcode.Convert)
		return
	}

	transaction, err := a.DB.Begin(ctx)
	if err != nil {
		Err(w, errcode.Database)
		return
	}
	defer func() {
		err := transaction.Rollback(ctx)
		if !errors.Is(err, pgx.ErrTxClosed) {
			panic(err)
		}
	}()

	query := model.New(transaction)

	exist, err := query.CheckUserByUsername(ctx, params.Username)
	if err != nil {
		Err(w, errcode.Database)
		return
	}
	if exist {
		Err(w, errcode.UsernameOccupy)
		return
	}

	user, err := query.CreateUser(ctx, params)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	var userRoleList []UserRole
	for _, req := range req.Role {
		params, err := createUserRoleParams(req, user.ID)
		if err != nil {
			Err(w, errcode.Convert)
			return
		}

		userRole, err := query.CreateUserRole(ctx, params)
		if err != nil {
			Err(w, errcode.Database)
			return
		}

		userRoleList = append(userRoleList, userRoleResp(userRole))
	}

	err = transaction.Commit(ctx)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	resp := userResp(user)
	resp.Role = userRoleList
	encode(w, resp)
}

func (a *API) GetApiV1Users(w http.ResponseWriter, r *http.Request, params GetApiV1UsersParams) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	err := validator.New().Struct(params)
	if err != nil {
		Err(w, errcode.Validate)
		return
	}

	var modelParams model.ListUserParams
	modelParams.Column1 = params.Username
	modelParams.Column2 = params.Name
	modelParams.Column3 = params.Status
	modelParams.ID, modelParams.Limit = paging(params.Current, params.PageSize)

	query := model.New(a.DB)

	userList, err := query.ListUser(ctx, modelParams)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	var userIDList []int32
	var respList []User
	for _, user := range userList {
		userIDList = append(userIDList, user.ID)
		respList = append(respList, userResp(user))
	}

	userRoleList, err := query.ListUserRoleByUserIDList(ctx, userIDList)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	userIDToUserRoleList := make(map[int32][]UserRole)
	for _, userRole := range userRoleList {
		userIDToUserRoleList[userRole.UserID] = append(userIDToUserRoleList[userRole.UserID], userRoleResp(userRole))
	}

	for i := range respList {
		respList[i].Role = userIDToUserRoleList[*respList[i].Id]
	}

	encode(w, respList)
}

func (a *API) DeleteApiV1UsersId(w http.ResponseWriter, r *http.Request, id int32) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	transaction, err := a.DB.Begin(ctx)
	if err != nil {
		Err(w, errcode.Database)
		return
	}
	defer func() {
		err := transaction.Rollback(ctx)
		if !errors.Is(err, pgx.ErrTxClosed) {
			panic(err)
		}
	}()

	query := model.New(transaction)

	exist, err := query.CheckUserByID(ctx, id)
	if err != nil {
		Err(w, errcode.Database)
		return
	}
	if !exist {
		Err(w, errcode.UserNotExist)
		return
	}

	err = query.DeleteUser(ctx, id)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	err = query.DeleteUserRoleByUserID(ctx, id)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	err = transaction.Commit(ctx)
	if err != nil {
		Err(w, errcode.Database)
		return
	}
}

func (a *API) GetApiV1UsersId(w http.ResponseWriter, r *http.Request, id int32) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	query := model.New(a.DB)

	user, err := query.GetUser(ctx, id)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		Err(w, errcode.Database)
		return
	}
	if errors.Is(err, pgx.ErrNoRows) {
		Err(w, errcode.UserNotExist)
		return
	}

	userRoleList, err := query.ListUserRoleByUserIDList(ctx, []int32{user.ID})
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	var userRoleListResp []UserRole
	for _, userRole := range userRoleList {
		userRoleListResp = append(userRoleListResp, userRoleResp(userRole))
	}

	resp := userResp(user)
	resp.Role = userRoleListResp
	encode(w, resp)
}

func (a *API) PutApiV1UsersId(w http.ResponseWriter, r *http.Request, id int32) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	var req User
	decode(w, r, &req)
	err := validator.New().Struct(req)
	if err != nil {
		Err(w, errcode.Validate)
		return
	}

	params, err := updateUserParams(req)
	if err != nil {
		Err(w, errcode.Convert)
		return
	}

	transaction, err := a.DB.Begin(ctx)
	if err != nil {
		Err(w, errcode.Database)
		return
	}
	defer func() {
		err := transaction.Rollback(ctx)
		if !errors.Is(err, pgx.ErrTxClosed) {
			panic(err)
		}
	}()

	query := model.New(transaction)

	userByGet, err := query.GetUser(ctx, id)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		Err(w, errcode.Database)
		return
	}
	if errors.Is(err, pgx.ErrNoRows) {
		Err(w, errcode.UserNotExist)
		return
	}

	// update username
	if req.Username != userByGet.Username {
		exist, err := query.CheckUserByUsername(ctx, req.Username)
		if err != nil {
			Err(w, errcode.Database)
			return
		}
		if exist {
			Err(w, errcode.UsernameOccupy)
			return
		}
	}

	params.ID = id

	userByUpdate, err := query.UpdateUser(ctx, params)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	err = query.DeleteUserRoleByUserID(ctx, userByUpdate.ID)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	var userRoleList []UserRole
	for _, req := range req.Role {
		params, err := createUserRoleParams(req, userByUpdate.ID)
		if err != nil {
			Err(w, errcode.Convert)
			return
		}

		userRole, err := query.CreateUserRole(ctx, params)
		if err != nil {
			Err(w, errcode.Database)
			return
		}
		
		userRoleList = append(userRoleList, userRoleResp(userRole))
	}

	err = transaction.Commit(ctx)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	resp := userResp(userByUpdate)
	resp.Role = userRoleList
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

func userResp(m model.AppUser) User {
	var resp User
	resp.Id = &m.ID
	resp.Username = m.Username
	resp.Name = m.Name
	resp.Email = m.Email
	resp.Phone = m.Phone
	resp.Remark = m.Remark
	resp.Status = UserStatus(m.Status)
	resp.Created = m.Created.Time.Format(pgTimestampFormat)
	resp.Updated = m.Updated.Time.Format(pgTimestampFormat)
	return resp
}

func createUserRoleParams(req UserRole, userID int32) (model.CreateUserRoleParams, error) {
	var params model.CreateUserRoleParams
	params.UserID = userID
	params.RoleID = req.RoleId
	if req.Created == "" {
		req.Created = time.Now().Format(pgTimestampFormat)
	}
	err := params.Created.Scan(req.Created)
	if err != nil {
		return model.CreateUserRoleParams{}, err
	}
	if req.Updated == "" {
		req.Updated = time.Now().Format(pgTimestampFormat)
	}
	err = params.Updated.Scan(req.Updated)
	if err != nil {
		return model.CreateUserRoleParams{}, err
	}
	return params, nil
}

func userRoleResp(m model.UserRole) UserRole {
	var resp UserRole
	resp.Id = &m.ID
	resp.UserId = &m.UserID
	resp.RoleId = m.RoleID
	resp.Created = m.Created.Time.Format(pgTimestampFormat)
	resp.Updated = m.Updated.Time.Format(pgTimestampFormat)
	return resp
}

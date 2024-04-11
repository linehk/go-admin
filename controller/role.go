package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/linehk/go-admin/errcode"
	"github.com/linehk/go-admin/model"
)

func (a *API) PostApiV1Roles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	var req Role
	decode(w, r, &req)
	err := validator.New().Struct(req)
	if err != nil {
		Err(w, errcode.Validate)
		return
	}

	params, err := createRoleParams(req)
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

	exist, err := query.CheckRoleByCode(ctx, params.Code)
	if err != nil {
		Err(w, errcode.Database)
		return
	}
	if exist {
		Err(w, errcode.CodeOccupy)
		return
	}

	role, err := query.CreateRole(ctx, params)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	var roleMenuList []RoleMenu
	for _, req := range req.Menu {
		params, err := createRoleMenuParams(req, role.ID)
		if err != nil {
			Err(w, errcode.Convert)
			return
		}

		roleMenu, err := query.CreateRoleMenu(ctx, params)
		if err != nil {
			Err(w, errcode.Database)
			return
		}

		roleMenuList = append(roleMenuList, roleMenuResp(roleMenu))
	}

	err = transaction.Commit(ctx)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	resp := roleResp(role)
	resp.Menu = roleMenuList
	encode(w, resp)
}

func (a *API) GetApiV1Roles(w http.ResponseWriter, r *http.Request, params GetApiV1RolesParams) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	err := validator.New().Struct(params)
	if err != nil {
		Err(w, errcode.Validate)
		return
	}

	var modelParams model.ListRoleParams
	modelParams.Column1 = params.Name
	modelParams.Column2 = params.Status
	modelParams.ID, modelParams.Limit = paging(params.Current, params.PageSize)

	query := model.New(a.DB)

	roleList, err := query.ListRole(ctx, modelParams)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	var roleIDList []int32
	var respList []Role
	for _, role := range roleList {
		roleIDList = append(roleIDList, role.ID)
		respList = append(respList, roleResp(role))
	}

	roleMenuList, err := query.ListRoleMenuByRoleIDList(ctx, roleIDList)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	roleIDToRoleMenuList := make(map[int32][]RoleMenu)
	for _, roleMenu := range roleMenuList {
		roleIDToRoleMenuList[roleMenu.RoleID] = append(roleIDToRoleMenuList[roleMenu.RoleID], roleMenuResp(roleMenu))
	}

	for i := range respList {
		respList[i].Menu = roleIDToRoleMenuList[*respList[i].Id]
	}

	encode(w, respList)
}

func (a *API) DeleteApiV1RolesId(w http.ResponseWriter, r *http.Request, id int32) {
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

	exist, err := query.CheckRoleByID(ctx, id)
	if err != nil {
		Err(w, errcode.Database)
		return
	}
	if !exist {
		Err(w, errcode.RoleNotExist)
		return
	}

	err = query.DeleteRole(ctx, id)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	err = query.DeleteRoleMenuByRoleID(ctx, id)
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

func (a *API) GetApiV1RolesId(w http.ResponseWriter, r *http.Request, id int32) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	query := model.New(a.DB)

	role, err := query.GetRole(ctx, id)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		Err(w, errcode.Database)
		return
	}
	if errors.Is(err, pgx.ErrNoRows) {
		Err(w, errcode.RoleNotExist)
		return
	}

	roleMenuList, err := query.ListRoleMenuByRoleIDList(ctx, []int32{role.ID})
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	var roleMenuListResp []RoleMenu
	for _, roleMenu := range roleMenuList {
		roleMenuListResp = append(roleMenuListResp, roleMenuResp(roleMenu))
	}

	resp := roleResp(role)
	resp.Menu = roleMenuListResp
	encode(w, resp)
}

func (a *API) PutApiV1RolesId(w http.ResponseWriter, r *http.Request, id int32) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	var req Role
	decode(w, r, &req)
	err := validator.New().Struct(req)
	if err != nil {
		Err(w, errcode.Validate)
		return
	}

	params, err := updateRoleParams(req)
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

	roleByGet, err := query.GetRole(ctx, id)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		Err(w, errcode.Database)
		return
	}
	if errors.Is(err, pgx.ErrNoRows) {
		Err(w, errcode.RoleNotExist)
		return
	}

	// update code
	if req.Code != roleByGet.Code {
		exist, err := query.CheckRoleByCode(ctx, req.Code)
		if err != nil {
			Err(w, errcode.Database)
			return
		}
		if exist {
			Err(w, errcode.CodeOccupy)
			return
		}
	}

	params.ID = id

	roleByUpdate, err := query.UpdateRole(ctx, params)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	err = query.DeleteRoleMenuByRoleID(ctx, roleByUpdate.ID)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	var roleMenuList []RoleMenu
	for _, req := range req.Menu {
		params, err := createRoleMenuParams(req, roleByUpdate.ID)
		if err != nil {
			Err(w, errcode.Convert)
			return
		}

		roleMenu, err := query.CreateRoleMenu(ctx, params)
		if err != nil {
			Err(w, errcode.Database)
			return
		}
		
		roleMenuList = append(roleMenuList, roleMenuResp(roleMenu))
	}

	err = transaction.Commit(ctx)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	resp := roleResp(roleByUpdate)
	resp.Menu = roleMenuList
	encode(w, resp)
}

func createRoleParams(req Role) (model.CreateRoleParams, error) {
	var params model.CreateRoleParams
	params.Code = req.Code
	params.Name = req.Name
	params.Description = req.Description
	params.Sequence = req.Sequence
	if req.Status == "" {
		req.Status = Enabled
	}
	params.Status = string(req.Status)
	if req.Created == "" {
		req.Created = time.Now().Format(pgTimestampFormat)
	}
	err := params.Created.Scan(req.Created)
	if err != nil {
		return model.CreateRoleParams{}, err
	}
	if req.Updated == "" {
		req.Updated = time.Now().Format(pgTimestampFormat)
	}
	err = params.Updated.Scan(req.Updated)
	if err != nil {
		return model.CreateRoleParams{}, err
	}
	return params, nil
}

func updateRoleParams(req Role) (model.UpdateRoleParams, error) {
	var params model.UpdateRoleParams
	params.Code = req.Code
	params.Name = req.Name
	params.Description = req.Description
	params.Sequence = req.Sequence
	params.Status = string(req.Status)
	err := params.Created.Scan(req.Created)
	if err != nil {
		return model.UpdateRoleParams{}, err
	}
	if req.Updated == "" {
		req.Updated = time.Now().Format(pgTimestampFormat)
	}
	err = params.Updated.Scan(req.Updated)
	if err != nil {
		return model.UpdateRoleParams{}, err
	}
	return params, nil
}

func roleResp(roleModel model.Role) Role {
	var resp Role
	resp.Id = &roleModel.ID
	resp.Code = roleModel.Code
	resp.Name = roleModel.Name
	resp.Description = roleModel.Description
	resp.Sequence = roleModel.Sequence
	resp.Status = RoleStatus(roleModel.Status)
	resp.Created = roleModel.Created.Time.Format(pgTimestampFormat)
	resp.Updated = roleModel.Updated.Time.Format(pgTimestampFormat)
	return resp
}

func createRoleMenuParams(req RoleMenu, roleID int32) (model.CreateRoleMenuParams, error) {
	var params model.CreateRoleMenuParams
	params.RoleID = roleID
	params.MenuID = req.MenuId
	if req.Created == "" {
		req.Created = time.Now().Format(pgTimestampFormat)
	}
	err := params.Created.Scan(req.Created)
	if err != nil {
		return model.CreateRoleMenuParams{}, err
	}
	if req.Updated == "" {
		req.Updated = time.Now().Format(pgTimestampFormat)
	}
	err = params.Updated.Scan(req.Updated)
	if err != nil {
		return model.CreateRoleMenuParams{}, err
	}
	return params, nil
}

func roleMenuResp(m model.RoleMenu) RoleMenu {
	var resp RoleMenu
	resp.Id = &m.ID
	resp.RoleId = &m.RoleID
	resp.MenuId = m.MenuID
	resp.Created = m.Created.Time.Format(pgTimestampFormat)
	resp.Updated = m.Updated.Time.Format(pgTimestampFormat)
	return resp
}

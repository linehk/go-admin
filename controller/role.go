package controller

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/linehk/go-admin/errcode"
	"github.com/linehk/go-admin/model"
)

func (a *API) PostApiV1Roles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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

	role, err := a.DB.CreateRole(r.Context(), params)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	resp := roleResp(role)
	encode(w, resp)
}

func (a *API) GetApiV1Roles(w http.ResponseWriter, r *http.Request, params GetApiV1RolesParams) {
	w.Header().Set("Content-Type", "application/json")
	
	err := validator.New().Struct(params)
	if err != nil {
		Err(w, errcode.Validate)
		return
	}

	var modelParams model.ListRoleParams
	modelParams.Column1 = params.Name
	modelParams.Column2 = params.Status
	modelParams.ID, modelParams.Limit = paging(params.Current, params.PageSize)

	roleList, err := a.DB.ListRole(r.Context(), modelParams)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	var respList []Role
	for _, roleModel := range roleList {
		respList = append(respList, roleResp(roleModel))
	}

	encode(w, respList)
}

func (a *API) DeleteApiV1RolesId(w http.ResponseWriter, r *http.Request, id int32) {
	w.Header().Set("Content-Type", "application/json")

	err := a.DB.DeleteRole(r.Context(), id)
	if err != nil {
		Err(w, errcode.Database)
		return
	}
}

func (a *API) GetApiV1RolesId(w http.ResponseWriter, r *http.Request, id int32) {
	w.Header().Set("Content-Type", "application/json")

	role, err := a.DB.GetRole(r.Context(), id)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	resp := roleResp(role)
	encode(w, resp)
}

func (a *API) PutApiV1RolesId(w http.ResponseWriter, r *http.Request, id int32) {
	w.Header().Set("Content-Type", "application/json")

	var req Role
	decode(w, r, &req)

	params, err := updateRoleParams(req)
	if err != nil {
		Err(w, errcode.Convert)
		return
	}

	params.ID = id
	role, err := a.DB.UpdateRole(r.Context(), params)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	resp := roleResp(role)
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

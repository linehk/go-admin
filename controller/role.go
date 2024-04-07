package controller

import (
	"net/http"
	"time"

	"github.com/linehk/go-admin/errcode"
	"github.com/linehk/go-admin/model"
)

func (a *API) PostApiV1Roles(w http.ResponseWriter, r *http.Request) {
	var req Role
	decode(w, r, &req)

	createRoleParams, err := reqToCreateRoleParams(req)
	if err != nil {
		ReturnErr(w, errcode.Convert)
	}

	roleModel, err := a.DB.CreateRole(r.Context(), createRoleParams)
	if err != nil {
		ReturnErr(w, errcode.Database)
	}

	roleResp := roleModelToResp(roleModel)
	encode(w, roleResp)
}

func (a *API) GetApiV1Roles(w http.ResponseWriter, r *http.Request, params GetApiV1RolesParams) {
	w.Header().Set("Content-Type", "application/json")
	var listRoleParams model.ListRoleParams
	if params.Name != nil {
		listRoleParams.Column1 = *params.Name
	}
	if params.Status != nil {
		listRoleParams.Column2 = *params.Status
	}

	listRoleParams.ID, listRoleParams.Limit = paging(params.Current, params.PageSize)

	roleModelList, err := a.DB.ListRole(r.Context(), listRoleParams)
	if err != nil {
		ReturnErr(w, errcode.Database)
	}

	var roleRespList []Role
	for _, roleModel := range roleModelList {
		roleRespList = append(roleRespList, roleModelToResp(roleModel))
	}

	encode(w, roleRespList)
}

func (a *API) DeleteApiV1RolesId(w http.ResponseWriter, r *http.Request, id int32) {
	w.Header().Set("Content-Type", "application/json")
	err := a.DB.DeleteRole(r.Context(), id)
	if err != nil {
		ReturnErr(w, errcode.Database)
	}
}

func (a *API) GetApiV1RolesId(w http.ResponseWriter, r *http.Request, id int32) {
	w.Header().Set("Content-Type", "application/json")
	roleModel, err := a.DB.GetRole(r.Context(), id)
	if err != nil {
		ReturnErr(w, errcode.Database)
	}

	roleResp := roleModelToResp(roleModel)
	encode(w, roleResp)
}

func (a *API) PutApiV1RolesId(w http.ResponseWriter, r *http.Request, id int32) {
	var req Role
	decode(w, r, &req)

	updateRoleParams, err := reqToUpdateRoleParams(req)
	if err != nil {
		ReturnErr(w, errcode.Convert)
	}

	updateRoleParams.ID = id
	roleModel, err := a.DB.UpdateRole(r.Context(), updateRoleParams)
	if err != nil {
		ReturnErr(w, errcode.Database)
	}

	roleResp := roleModelToResp(roleModel)
	encode(w, roleResp)
}

func reqToCreateRoleParams(req Role) (model.CreateRoleParams, error) {
	var createRoleParams model.CreateRoleParams
	createRoleParams.Code = req.Code
	createRoleParams.Name = req.Name
	if req.Description != nil {
		createRoleParams.Description = *req.Description
	}
	if req.Sequence != nil {
		createRoleParams.Sequence = int16(*req.Sequence)
	}
	if req.Status != nil {
		createRoleParams.Status = string(*req.Status)
	} else {
		createRoleParams.Status = string(Enabled)
	}
	if req.Created != nil {
		err := createRoleParams.Created.Scan(*req.Created)
		if err != nil {
			return model.CreateRoleParams{}, err
		}
	} else {
		err := createRoleParams.Created.Scan(time.Now())
		if err != nil {
			return model.CreateRoleParams{}, err
		}
	}
	if req.Updated != nil {
		err := createRoleParams.Updated.Scan(*req.Updated)
		if err != nil {
			return model.CreateRoleParams{}, err
		}
	} else {
		err := createRoleParams.Updated.Scan(time.Now())
		if err != nil {
			return model.CreateRoleParams{}, err
		}
	}
	return createRoleParams, nil
}

func reqToUpdateRoleParams(req Role) (model.UpdateRoleParams, error) {
	var updateRoleParams model.UpdateRoleParams
	updateRoleParams.Code = req.Code
	updateRoleParams.Name = req.Name
	if req.Description != nil {
		updateRoleParams.Description = *req.Description
	}
	if req.Sequence != nil {
		updateRoleParams.Sequence = int16(*req.Sequence)
	}
	if req.Status != nil {
		updateRoleParams.Status = string(*req.Status)
	} else {
		updateRoleParams.Status = string(Enabled)
	}
	if req.Created != nil {
		err := updateRoleParams.Created.Scan(*req.Created)
		if err != nil {
			return model.UpdateRoleParams{}, err
		}
	}
	if req.Updated != nil {
		err := updateRoleParams.Updated.Scan(*req.Updated)
		if err != nil {
			return model.UpdateRoleParams{}, err
		}
	} else {
		err := updateRoleParams.Updated.Scan(time.Now())
		if err != nil {
			return model.UpdateRoleParams{}, err
		}
	}
	return updateRoleParams, nil
}

func roleModelToResp(roleModel model.Role) Role {
	const format = "2006-01-02 15:04:05.999999999"
	var roleResp Role
	roleResp.Code = roleModel.Code
	roleResp.Name = roleModel.Name
	roleResp.Description = &roleModel.Description
	sequence := int(roleModel.Sequence)
	roleResp.Sequence = &sequence
	roleStatus := RoleStatus(roleModel.Status)
	roleResp.Status = &roleStatus
	createdStr := roleModel.Created.Time.Format(format)
	roleResp.Created = &createdStr
	updatedStr := roleModel.Updated.Time.Format(format)
	roleResp.Updated = &updatedStr
	return roleResp
}

package controller

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/linehk/go-admin/errcode"
	"github.com/linehk/go-admin/model"
)

func (a *API) PostApiV1Menus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	var req Menu
	decode(w, r, &req)
	err := validator.New().Struct(req)
	if err != nil {
		Err(w, errcode.Validate)
		return
	}

	params, err := createMenuParams(req)
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

	// have parent
	if req.ParentId != 0 {
		parentMenu, err := query.GetMenu(ctx, req.ParentId)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			Err(w, errcode.Database)
			return
		}

		if errors.Is(err, pgx.ErrNoRows) {
			Err(w, errcode.MenuNotExist)
			return
		}

		params.ParentPath = parentMenu.ParentPath + strconv.Itoa(int(parentMenu.ID)) + "."
	}

	checkParams := model.CheckMenuByCodeAndParentIDParams{
		Code:     params.Code,
		ParentID: params.ParentID,
	}
	exist, err := query.CheckMenuByCodeAndParentID(ctx, checkParams)
	if err != nil {
		Err(w, errcode.Database)
		return
	}
	if exist {
		Err(w, errcode.MenuCodeOccupy)
		return
	}

	menu, err := query.CreateMenu(ctx, params)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	var resourceList []Resource
	for _, req := range req.Resource {
		params, err := createResourceParams(req, menu.ID)
		if err != nil {
			Err(w, errcode.Convert)
			return
		}

		resource, err := query.CreateResource(ctx, params)
		if err != nil {
			Err(w, errcode.Database)
			return
		}

		resourceList = append(resourceList, resourceResp(resource))
	}

	err = transaction.Commit(ctx)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	resp := menuResp(menu)
	resp.Resource = resourceList
	encode(w, resp)
}

func (a *API) GetApiV1Menus(w http.ResponseWriter, r *http.Request, params GetApiV1MenusParams) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	err := validator.New().Struct(params)
	if err != nil {
		Err(w, errcode.Validate)
		return
	}

	var modelParams model.ListMenuParams
	modelParams.Column1 = params.Menuname
	modelParams.Column2 = params.Name
	modelParams.Column3 = params.Status
	modelParams.ID, modelParams.Limit = paging(params.Current, params.PageSize)

	query := model.New(a.DB)

	menuList, err := query.ListMenu(ctx, modelParams)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	var menuIDList []int32
	var respList []Menu
	for _, menu := range menuList {
		menuIDList = append(menuIDList, menu.ID)
		respList = append(respList, menuResp(menu))
	}

	menuRoleList, err := query.ListMenuRoleByMenuIDList(ctx, menuIDList)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	menuIDToMenuRoleList := make(map[int32][]MenuRole)
	for _, menuRole := range menuRoleList {
		menuIDToMenuRoleList[menuRole.MenuID] = append(menuIDToMenuRoleList[menuRole.MenuID], resourceResp(menuRole))
	}

	for i := range respList {
		respList[i].Role = menuIDToMenuRoleList[*respList[i].Id]
	}

	encode(w, respList)
}

func (a *API) DeleteApiV1MenusId(w http.ResponseWriter, r *http.Request, id int32) {
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

	menu, err := query.GetMenu(ctx, id)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		Err(w, errcode.Database)
		return
	}
	if errors.Is(err, pgx.ErrNoRows) {
		Err(w, errcode.MenuNotExist)
		return
	}

	childIdList, err := query.ListChildID(ctx, menu.ParentPath+strconv.Itoa(int(menu.ID))+".")
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	err = query.DeleteMenu(ctx, id)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	err = query.DeleteMenuRoleByMenuID(ctx, id)
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

func (a *API) GetApiV1MenusId(w http.ResponseWriter, r *http.Request, id int32) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	query := model.New(a.DB)

	menu, err := query.GetMenu(ctx, id)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		Err(w, errcode.Database)
		return
	}
	if errors.Is(err, pgx.ErrNoRows) {
		Err(w, errcode.MenuNotExist)
		return
	}

	menuRoleList, err := query.ListMenuRoleByMenuIDList(ctx, []int32{menu.ID})
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	var menuRoleListResp []MenuRole
	for _, menuRole := range menuRoleList {
		menuRoleListResp = append(menuRoleListResp, resourceResp(menuRole))
	}

	resp := menuResp(menu)
	resp.Role = menuRoleListResp
	encode(w, resp)
}

func (a *API) PutApiV1MenusId(w http.ResponseWriter, r *http.Request, id int32) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	var req Menu
	decode(w, r, &req)
	err := validator.New().Struct(req)
	if err != nil {
		Err(w, errcode.Validate)
		return
	}

	params, err := updateMenuParams(req)
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

	menuByGet, err := query.GetMenu(ctx, id)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		Err(w, errcode.Database)
		return
	}
	if errors.Is(err, pgx.ErrNoRows) {
		Err(w, errcode.MenuNotExist)
		return
	}

	// update menuname
	if req.Menuname != menuByGet.Menuname {
		exist, err := query.CheckMenuByMenuname(ctx, req.Menuname)
		if err != nil {
			Err(w, errcode.Database)
			return
		}
		if exist {
			Err(w, errcode.MenunameOccupy)
			return
		}
	}

	params.ID = id

	menuByUpdate, err := query.UpdateMenu(ctx, params)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	err = query.DeleteMenuRoleByMenuID(ctx, menuByUpdate.ID)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	var menuRoleList []MenuRole
	for _, req := range req.Role {
		params, err := createResourceParams(req, menuByUpdate.ID)
		if err != nil {
			Err(w, errcode.Convert)
			return
		}

		menuRole, err := query.CreateMenuRole(ctx, params)
		if err != nil {
			Err(w, errcode.Database)
			return
		}

		menuRoleList = append(menuRoleList, resourceResp(menuRole))
	}

	err = transaction.Commit(ctx)
	if err != nil {
		Err(w, errcode.Database)
		return
	}

	resp := menuResp(menuByUpdate)
	resp.Role = menuRoleList
	encode(w, resp)
}

func createMenuParams(req Menu) (model.CreateMenuParams, error) {
	var params model.CreateMenuParams
	params.Code = req.Code
	params.Name = req.Name
	params.Description = req.Description
	params.Sequence = req.Sequence
	if req.Type == "" {
		req.Type = Button
	}
	params.Type = string(req.Type)
	params.Path = req.Path
	params.Property = req.Property
	params.ParentID = req.ParentId
	params.ParentPath = req.ParentPath
	if req.Status == "" {
		req.Status = MenuStatusEnabled
	}
	params.Status = string(req.Status)
	if req.Created == "" {
		req.Created = time.Now().Format(pgTimestampFormat)
	}
	err := params.Created.Scan(req.Created)
	if err != nil {
		return model.CreateMenuParams{}, err
	}
	if req.Updated == "" {
		req.Updated = time.Now().Format(pgTimestampFormat)
	}
	err = params.Updated.Scan(req.Updated)
	if err != nil {
		return model.CreateMenuParams{}, err
	}
	return params, nil
}

func updateMenuParams(req Menu) (model.UpdateMenuParams, error) {
	var params model.UpdateMenuParams
	params.Menuname = req.Menuname
	password, err := hash(req.Password)
	if err != nil {
		return model.UpdateMenuParams{}, err
	}
	params.Password = password
	params.Name = req.Name
	params.Email = req.Email
	params.Phone = req.Phone
	params.Remark = req.Remark
	params.Status = string(req.Status)
	err = params.Created.Scan(req.Created)
	if err != nil {
		return model.UpdateMenuParams{}, err
	}
	if req.Updated == "" {
		req.Updated = time.Now().Format(pgTimestampFormat)
	}
	err = params.Updated.Scan(req.Updated)
	if err != nil {
		return model.UpdateMenuParams{}, err
	}
	return params, nil
}

func menuResp(m model.Menu) Menu {
	var resp Menu
	resp.Id = &m.ID
	resp.Code = m.Code
	resp.Name = m.Name
	resp.Description = m.Description
	resp.Sequence = m.Sequence
	resp.Type = MenuType(m.Type)
	resp.Path = m.Path
	resp.Property = m.Property
	resp.ParentId = m.ParentID
	resp.ParentPath = m.ParentPath
	resp.Status = MenuStatus(m.Status)
	resp.Created = m.Created.Time.Format(pgTimestampFormat)
	resp.Updated = m.Updated.Time.Format(pgTimestampFormat)
	return resp
}

func createResourceParams(req Resource, menuID int32) (model.CreateResourceParams, error) {
	var params model.CreateResourceParams
	params.MenuID = menuID
	params.Method = req.Method
	params.Path = req.Path
	if req.Created == "" {
		req.Created = time.Now().Format(pgTimestampFormat)
	}
	err := params.Created.Scan(req.Created)
	if err != nil {
		return model.CreateResourceParams{}, err
	}
	if req.Updated == "" {
		req.Updated = time.Now().Format(pgTimestampFormat)
	}
	err = params.Updated.Scan(req.Updated)
	if err != nil {
		return model.CreateResourceParams{}, err
	}
	return params, nil
}

func resourceResp(m model.Resource) Resource {
	var resp Resource
	resp.Id = &m.ID
	resp.MenuId = m.MenuID
	resp.Method = m.Method
	resp.Path = m.Path
	resp.Created = m.Created.Time.Format(pgTimestampFormat)
	resp.Updated = m.Updated.Time.Format(pgTimestampFormat)
	return resp
}

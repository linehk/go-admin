package role

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/linehk/go-admin/controller"
	"github.com/linehk/go-admin/tests"
	"github.com/stretchr/testify/assert"
)

var (
	role1JSON = `{
"code": "code1",
"name": "name1",
"description": "description1",
"sequence": 1,
"status": "enabled",
"created": "2024-04-04 13:56:35.671521",
"updated": "2024-04-05 13:56:35.671521",
"menu": [
    {
        "menu_id": 1,
        "created": "2024-04-04 13:56:35.671521",
        "updated": "2024-04-05 13:56:35.671521"
    },
    {
        "menu_id": 2,
        "created": "2024-04-04 13:56:35.671521",
        "updated": "2024-04-05 13:56:35.671521"
    }
]
}`

	id1         int32 = 1
	roleMenuId1 int32 = 1
	roleMenuId2 int32 = 2
	role1             = controller.Role{
		Id:          &id1,
		Code:        "code1",
		Name:        "name1",
		Description: "description1",
		Sequence:    1,
		Status:      controller.RoleStatusEnabled,
		Created:     "2024-04-04 13:56:35.671521",
		Updated:     "2024-04-05 13:56:35.671521",
		Menu: []controller.RoleMenu{
			{
				Id:      &roleMenuId1,
				RoleId:  &id1,
				MenuId:  1,
				Created: "2024-04-04 13:56:35.671521",
				Updated: "2024-04-05 13:56:35.671521",
			},
			{
				Id:      &roleMenuId2,
				RoleId:  &id1,
				MenuId:  2,
				Created: "2024-04-04 13:56:35.671521",
				Updated: "2024-04-05 13:56:35.671521",
			},
		},
	}
)

var (
	role2JSON = `{
"code": "code2",
"name": "name2",
"description": "description2",
"sequence": 2,
"status": "enabled",
"created": "2024-03-04 13:56:35.671521",
"updated": "2024-03-05 13:56:35.671521",
"menu": [
    {
        "menu_id": 3,
        "created": "2024-04-04 13:56:35.671521",
        "updated": "2024-04-05 13:56:35.671521"
    },
    {
        "menu_id": 4,
        "created": "2024-04-04 13:56:35.671521",
        "updated": "2024-04-05 13:56:35.671521"
    }
]
}`

	id2         int32 = 2
	roleMenuId3 int32 = 3
	roleMenuId4 int32 = 4
	role2             = controller.Role{
		Id:          &id2,
		Code:        "code2",
		Name:        "name2",
		Description: "description2",
		Sequence:    2,
		Status:      controller.RoleStatusEnabled,
		Created:     "2024-03-04 13:56:35.671521",
		Updated:     "2024-03-05 13:56:35.671521",
		Menu: []controller.RoleMenu{
			{
				Id:      &roleMenuId3,
				RoleId:  &id2,
				MenuId:  3,
				Created: "2024-04-04 13:56:35.671521",
				Updated: "2024-04-05 13:56:35.671521",
			},
			{
				Id:      &roleMenuId4,
				RoleId:  &id2,
				MenuId:  4,
				Created: "2024-04-04 13:56:35.671521",
				Updated: "2024-04-05 13:56:35.671521",
			},
		},
	}
)

func createRole(db *pgx.Conn, reqJSON string) controller.Role {
	req := httptest.NewRequest(http.MethodPost, tests.BaseURL+"api/v1/roles", strings.NewReader(reqJSON))
	r := httptest.NewRecorder()
	api := &controller.API{DB: db}
	api.PostApiV1Roles(r, req)
	var actual controller.Role
	_ = json.NewDecoder(r.Body).Decode(&actual)
	return actual
}

func TestPostApiV1Roles(t *testing.T) {
	db := tests.ContainerDB(t)
	actual := createRole(db, role1JSON)

	assert.Equal(t, role1, actual)
}

func TestGetApiV1RolesId(t *testing.T) {
	db := tests.ContainerDB(t)
	_ = createRole(db, role1JSON)

	var id int32 = 1
	req := httptest.NewRequest(http.MethodGet, tests.BaseURL+fmt.Sprintf("api/v1/roles/%d", id), nil)
	r := httptest.NewRecorder()

	api := &controller.API{DB: db}
	api.GetApiV1RolesId(r, req, id)
	var actual controller.Role
	_ = json.NewDecoder(r.Body).Decode(&actual)

	assert.Equal(t, role1, actual)
}

func TestDeleteApiV1RolesId(t *testing.T) {
	db := tests.ContainerDB(t)
	_ = createRole(db, role1JSON)

	var id int32 = 1
	req := httptest.NewRequest(http.MethodDelete, tests.BaseURL+fmt.Sprintf("api/v1/roles/%d", id), nil)
	r := httptest.NewRecorder()

	api := &controller.API{DB: db}
	api.DeleteApiV1RolesId(r, req, id)

	assert.Equal(t, http.StatusOK, r.Code)
}

func TestPutApiV1RolesId(t *testing.T) {
	db := tests.ContainerDB(t)
	_ = createRole(db, role1JSON)

	var id int32 = 1
	req := httptest.NewRequest(http.MethodPut, tests.BaseURL+fmt.Sprintf("api/v1/roles/%d", id), strings.NewReader(role2JSON))
	r := httptest.NewRecorder()

	api := &controller.API{DB: db}
	api.PutApiV1RolesId(r, req, id)

	var actual controller.Role
	_ = json.NewDecoder(r.Body).Decode(&actual)
	// ignore id
	actual.Id = &id2
	for i := range actual.Menu {
		actual.Menu[i].RoleId = &id2
	}

	assert.Equal(t, role2, actual)
}

func TestGetApiV1Roles(t *testing.T) {
	db := tests.ContainerDB(t)
	_ = createRole(db, role1JSON)
	_ = createRole(db, role2JSON)

	api := &controller.API{DB: db}
	req := httptest.NewRequest(http.MethodGet, tests.BaseURL+"api/v1/roles", nil)
	params := controller.GetApiV1RolesParams{
		Name:     "",
		Status:   string(controller.RoleStatusEnabled),
		Current:  0,
		PageSize: 10,
	}
	r := httptest.NewRecorder()
	api.GetApiV1Roles(r, req, params)

	var actual []controller.Role
	_ = json.NewDecoder(r.Body).Decode(&actual)
	expected := []controller.Role{
		role1,
		role2,
	}

	assert.Equal(t, expected, actual)
}

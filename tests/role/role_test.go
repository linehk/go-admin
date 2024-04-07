package role

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/linehk/go-admin/controller"
	"github.com/linehk/go-admin/model"
	"github.com/linehk/go-admin/tests"
	"github.com/stretchr/testify/assert"
)

func createRole(db *model.Queries, createRoleReqJSON string) controller.Role {
	createRoleReq := httptest.NewRequest(http.MethodPost, tests.BaseURL+"api/v1/roles", strings.NewReader(createRoleReqJSON))
	createRoleRecorder := httptest.NewRecorder()
	api := &controller.API{DB: db}
	api.PostApiV1Roles(createRoleRecorder, createRoleReq)
	var actual controller.Role
	_ = json.NewDecoder(createRoleRecorder.Body).Decode(&actual)
	return actual
}

var role1JSON = `{
"name": "name1",
"code": "code1",
"description": "description1",
"sequence": 1,
"status": "enabled",
"created": "2024-04-04 13:56:35.671521",
"updated": "2024-04-05 13:56:35.671521"
}`

var name1 = "name1"
var code1 = "code1"
var description1 = "description1"
var sequence1 = 1
var status1 = controller.Enabled
var created1 = "2024-04-04 13:56:35.671521"
var updated1 = "2024-04-05 13:56:35.671521"
var role1 = controller.Role{
	Name:        name1,
	Code:        code1,
	Description: &description1,
	Sequence:    &sequence1,
	Status:      &status1,
	Created:     &created1,
	Updated:     &updated1,
}

var role2JSON = `{
"name": "name2",
"code": "code2",
"description": "description2",
"sequence": 2,
"status": "enabled",
"created": "2024-03-04 13:56:35.671521",
"updated": "2024-03-05 13:56:35.671521"
}`

var name2 = "name2"
var code2 = "code2"
var description2 = "description2"
var sequence2 = 2
var status2 = controller.Enabled
var created2 = "2024-03-04 13:56:35.671521"
var updated2 = "2024-03-05 13:56:35.671521"
var role2 = controller.Role{
	Name:        name2,
	Code:        code2,
	Description: &description2,
	Sequence:    &sequence2,
	Status:      &status2,
	Created:     &created2,
	Updated:     &updated2,
}

func TestPostApiV1Roles(t *testing.T) {
	db := tests.ContainerDB(t)
	actual := createRole(db, role1JSON)

	assert.Equal(t, role1, actual)
}

func TestGetApiV1RolesId(t *testing.T) {
	db := tests.ContainerDB(t)
	_ = createRole(db, role1JSON)

	var roleId int32 = 1
	getReq := httptest.NewRequest(http.MethodGet, tests.BaseURL+fmt.Sprintf("api/v1/roles/%d", roleId), nil)
	getRecorder := httptest.NewRecorder()

	api := &controller.API{DB: db}
	api.GetApiV1RolesId(getRecorder, getReq, roleId)
	var actual controller.Role
	_ = json.NewDecoder(getRecorder.Body).Decode(&actual)

	assert.Equal(t, role1, actual)
}

func TestDeleteApiV1RolesId(t *testing.T) {
	db := tests.ContainerDB(t)
	_ = createRole(db, role1JSON)

	var roleId int32 = 1
	deleteReq := httptest.NewRequest(http.MethodDelete, tests.BaseURL+fmt.Sprintf("api/v1/roles/%d", roleId), nil)
	deleteRecorder := httptest.NewRecorder()

	api := &controller.API{DB: db}
	api.DeleteApiV1RolesId(deleteRecorder, deleteReq, roleId)

	actual := deleteRecorder.Code
	expected := http.StatusOK

	assert.Equal(t, expected, actual)
}

func TestPutApiV1RolesId(t *testing.T) {
	db := tests.ContainerDB(t)
	_ = createRole(db, role1JSON)

	var roleId int32 = 1
	putReq := httptest.NewRequest(http.MethodPut, tests.BaseURL+fmt.Sprintf("api/v1/roles/%d", roleId), strings.NewReader(role2JSON))
	putRecorder := httptest.NewRecorder()

	api := &controller.API{DB: db}
	api.PutApiV1RolesId(putRecorder, putReq, roleId)

	var actual controller.Role
	_ = json.NewDecoder(putRecorder.Body).Decode(&actual)

	assert.Equal(t, role2, actual)
}

func TestGetApiV1Roles(t *testing.T) {
	db := tests.ContainerDB(t)
	_ = createRole(db, role1JSON)
	_ = createRole(db, role2JSON)

	api := &controller.API{DB: db}
	getReq := httptest.NewRequest(http.MethodGet, tests.BaseURL+"api/v1/roles", nil)
	status := string(controller.Enabled)
	params := controller.GetApiV1RolesParams{
		Current:  0,
		PageSize: 10,
		Status:   &status,
	}
	getRecorder := httptest.NewRecorder()
	api.GetApiV1Roles(getRecorder, getReq, params)

	var actual []controller.Role
	_ = json.NewDecoder(getRecorder.Body).Decode(&actual)
	expected := []controller.Role{
		role1,
		role2,
	}

	assert.Equal(t, expected, actual)
}

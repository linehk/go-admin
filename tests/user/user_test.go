package user

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

var (
	user1JSON = `{
"username": "username1",
"password": "password1",
"name": "name1",
"email": "example1@gmail.com",
"phone": "+14155552671",
"remark": "remark1",
"status": "activated",
"created": "2024-04-04 13:56:35.671521",
"updated": "2024-04-05 13:56:35.671521"
}`

	id1       int32 = 1
	username1       = "username1"
	name1           = "name1"
	email1          = "example1@gmail.com"
	phone1          = "+14155552671"
	remark1         = "remark1"
	status1         = controller.Activated
	created1        = "2024-04-04 13:56:35.671521"
	updated1        = "2024-04-05 13:56:35.671521"
	user1           = controller.User{
		Id:       &id1,
		Username: username1,
		Name:     name1,
		Email:    email1,
		Phone:    phone1,
		Remark:   remark1,
		Status:   status1,
		Created:  created1,
		Updated:  updated1,
	}
)

var (
	user2JSON = `{
"username": "username2",
"password": "password2",
"name": "name2",
"email": "example2@gmail.com",
"phone": "+442071838750",
"remark": "remark2",
"status": "activated",
"created": "2024-03-04 13:56:35.671521",
"updated": "2024-03-05 13:56:35.671521"
}`

	id2       int32 = 2
	username2       = "username2"
	name2           = "name2"
	email2          = "example2@gmail.com"
	phone2          = "+442071838750"
	remark2         = "remark2"
	status2         = controller.Activated
	created2        = "2024-03-04 13:56:35.671521"
	updated2        = "2024-03-05 13:56:35.671521"
	user2           = controller.User{
		Id:       &id2,
		Username: username2,
		Name:     name2,
		Email:    email2,
		Phone:    phone2,
		Remark:   remark2,
		Status:   status2,
		Created:  created2,
		Updated:  updated2,
	}
)

func createUser(db *model.Queries, reqJSON string) controller.User {
	req := httptest.NewRequest(http.MethodPost, tests.BaseURL+"api/v1/users", strings.NewReader(reqJSON))
	r := httptest.NewRecorder()
	api := &controller.API{DB: db}
	api.PostApiV1Users(r, req)
	var actual controller.User
	_ = json.NewDecoder(r.Body).Decode(&actual)
	return actual
}

func TestPostApiV1Users(t *testing.T) {
	db := tests.ContainerDB(t)
	actual := createUser(db, user1JSON)

	assert.Equal(t, user1, actual)
}

func TestGetApiV1UsersId(t *testing.T) {
	db := tests.ContainerDB(t)
	_ = createUser(db, user1JSON)

	var id int32 = 1
	req := httptest.NewRequest(http.MethodGet, tests.BaseURL+fmt.Sprintf("api/v1/users/%d", id), nil)
	r := httptest.NewRecorder()

	api := &controller.API{DB: db}
	api.GetApiV1UsersId(r, req, id)
	var actual controller.User
	_ = json.NewDecoder(r.Body).Decode(&actual)

	assert.Equal(t, user1, actual)
}

func TestDeleteApiV1UsersId(t *testing.T) {
	db := tests.ContainerDB(t)
	_ = createUser(db, user1JSON)

	var id int32 = 1
	req := httptest.NewRequest(http.MethodDelete, tests.BaseURL+fmt.Sprintf("api/v1/users/%d", id), nil)
	r := httptest.NewRecorder()

	api := &controller.API{DB: db}
	api.DeleteApiV1UsersId(r, req, id)

	assert.Equal(t, http.StatusOK, r.Code)
}

func TestPutApiV1UsersId(t *testing.T) {
	db := tests.ContainerDB(t)
	_ = createUser(db, user1JSON)

	var id int32 = 1
	req := httptest.NewRequest(http.MethodPut, tests.BaseURL+fmt.Sprintf("api/v1/users/%d", id), strings.NewReader(user2JSON))
	r := httptest.NewRecorder()

	api := &controller.API{DB: db}
	api.PutApiV1UsersId(r, req, id)

	var actual controller.User
	_ = json.NewDecoder(r.Body).Decode(&actual)
	actual.Id = &id2

	assert.Equal(t, user2, actual)
}

func TestGetApiV1Users(t *testing.T) {
	db := tests.ContainerDB(t)
	_ = createUser(db, user1JSON)
	_ = createUser(db, user2JSON)

	api := &controller.API{DB: db}
	req := httptest.NewRequest(http.MethodGet, tests.BaseURL+"api/v1/users", nil)
	params := controller.GetApiV1UsersParams{
		Username: "",
		Name:     "",
		Status:   string(controller.Activated),
		Current:  0,
		PageSize: 10,
	}
	r := httptest.NewRecorder()
	api.GetApiV1Users(r, req, params)

	var actual []controller.User
	_ = json.NewDecoder(r.Body).Decode(&actual)
	expected := []controller.User{
		user1,
		user2,
	}

	assert.Equal(t, expected, actual)
}

package user

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
	user1JSON = `{
"username": "username1",
"password": "password1",
"name": "name1",
"email": "example1@gmail.com",
"phone": "+14155552671",
"remark": "remark1",
"status": "activated",
"created": "2024-04-04 13:56:35.671521",
"updated": "2024-04-05 13:56:35.671521",
"role": [
    {
        "role_id": 1,
        "created": "2024-04-04 13:56:35.671521",
        "updated": "2024-04-05 13:56:35.671521"
    },
    {
        "role_id": 2,
        "created": "2024-04-04 13:56:35.671521",
        "updated": "2024-04-05 13:56:35.671521"
    }
]
}`

	id1         int32 = 1
	userRoleId1 int32 = 1
	userRoleId2 int32 = 2
	user1             = controller.User{
		Id:       &id1,
		Username: "username1",
		Name:     "name1",
		Email:    "example1@gmail.com",
		Phone:    "+14155552671",
		Remark:   "remark1",
		Status:   controller.Activated,
		Created:  "2024-04-04 13:56:35.671521",
		Updated:  "2024-04-05 13:56:35.671521",
		Role: []controller.UserRole{
			{
				Id:      &userRoleId1,
				UserId:  &id1,
				RoleId:  1,
				Created: "2024-04-04 13:56:35.671521",
				Updated: "2024-04-05 13:56:35.671521",
			},
			{
				Id:      &userRoleId2,
				UserId:  &id1,
				RoleId:  2,
				Created: "2024-04-04 13:56:35.671521",
				Updated: "2024-04-05 13:56:35.671521",
			},
		},
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
"updated": "2024-03-05 13:56:35.671521",
"role": [
    {
        "role_id": 3,
        "created": "2024-04-04 13:56:35.671521",
        "updated": "2024-04-05 13:56:35.671521"
    },
    {
        "role_id": 4,
        "created": "2024-04-04 13:56:35.671521",
        "updated": "2024-04-05 13:56:35.671521"
    }
]
}`

	id2         int32 = 2
	userRoleId3 int32 = 3
	userRoleId4 int32 = 4
	user2             = controller.User{
		Id:       &id2,
		Username: "username2",
		Name:     "name2",
		Email:    "example2@gmail.com",
		Phone:    "+442071838750",
		Remark:   "remark2",
		Status:   controller.Activated,
		Created:  "2024-03-04 13:56:35.671521",
		Updated:  "2024-03-05 13:56:35.671521",
		Role: []controller.UserRole{
			{
				Id:      &userRoleId3,
				UserId:  &id2,
				RoleId:  3,
				Created: "2024-04-04 13:56:35.671521",
				Updated: "2024-04-05 13:56:35.671521",
			},
			{
				Id:      &userRoleId4,
				UserId:  &id2,
				RoleId:  4,
				Created: "2024-04-04 13:56:35.671521",
				Updated: "2024-04-05 13:56:35.671521",
			},
		},
	}
)

func createUser(db *pgx.Conn, reqJSON string) controller.User {
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
	// ignore id
	actual.Id = &id2
	for i := range actual.Role {
		actual.Role[i].UserId = &id2
	}

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

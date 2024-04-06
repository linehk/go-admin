package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/linehk/go-admin/controller"
	"github.com/linehk/go-admin/model"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const baseURL = "http://localhost:8080/"

func ContainerDB(t *testing.T) *model.Queries {
	ctx := context.Background()
	pg, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16.2"),
		postgres.WithInitScripts(filepath.Join("..", "model", "schema.sql")),
		postgres.WithDatabase("go_admin"),
		postgres.WithUsername("dev"),
		postgres.WithPassword("dev"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)))
	if err != nil {
		t.Error(err)
	}
	dsn, err := pg.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Error(err)
	}
	return model.Setup(ctx, dsn)
}

func createUser(db *model.Queries, createUserReqJSON string) controller.User {
	createUserReq := httptest.NewRequest(http.MethodPost, baseURL+"api/v1/users", strings.NewReader(createUserReqJSON))
	createUserRecorder := httptest.NewRecorder()
	userImpl := &controller.UserImpl{DB: db}
	userImpl.PostApiV1Users(createUserRecorder, createUserReq)
	var actual controller.User
	_ = json.NewDecoder(createUserRecorder.Body).Decode(&actual)
	return actual
}

var user1JSON = `{
"username": "username1",
"password": "password1",
"email": "email1",
"phone": "phone1",
"remark": "remark1",
"status": "activated",
"created": "2024-04-04 13:56:35.671521",
"updated": "2024-04-05 13:56:35.671521"
}`

var username1 = "username1"
var name1 = ""
var email1 = "email1"
var phone1 = "phone1"
var remark1 = "remark1"
var status1 = controller.Activated
var created1 = "2024-04-04 13:56:35.671521"
var updated1 = "2024-04-05 13:56:35.671521"
var user1 = controller.User{
	Username: username1,
	Name:     &name1,
	Email:    &email1,
	Phone:    &phone1,
	Remark:   &remark1,
	Status:   &status1,
	Created:  &created1,
	Updated:  &updated1,
}

var user2JSON = `{
"username": "username2",
"password": "password2",
"email": "email2",
"phone": "phone2",
"remark": "remark2",
"status": "activated",
"created": "2024-03-04 13:56:35.671521",
"updated": "2024-03-05 13:56:35.671521"
}`

var username2 = "username2"
var name2 = ""
var email2 = "email2"
var phone2 = "phone2"
var remark2 = "remark2"
var status2 = controller.Activated
var created2 = "2024-03-04 13:56:35.671521"
var updated2 = "2024-03-05 13:56:35.671521"
var user2 = controller.User{
	Username: username2,
	Name:     &name2,
	Email:    &email2,
	Phone:    &phone2,
	Remark:   &remark2,
	Status:   &status2,
	Created:  &created2,
	Updated:  &updated2,
}

func TestPostApiV1Users(t *testing.T) {
	db := ContainerDB(t)
	actual := createUser(db, user1JSON)

	assert.Equal(t, user1, actual)
}

func TestGetApiV1UsersId(t *testing.T) {
	db := ContainerDB(t)
	_ = createUser(db, user1JSON)

	var userId int32 = 1
	getReq := httptest.NewRequest(http.MethodGet, baseURL+fmt.Sprintf("api/v1/users/%d", userId), nil)
	getRecorder := httptest.NewRecorder()

	userImpl := &controller.UserImpl{DB: db}
	userImpl.GetApiV1UsersId(getRecorder, getReq, userId)
	var actual controller.User
	_ = json.NewDecoder(getRecorder.Body).Decode(&actual)

	assert.Equal(t, user1, actual)
}

func TestDeleteApiV1UsersId(t *testing.T) {
	db := ContainerDB(t)
	_ = createUser(db, user1JSON)

	var userId int32 = 1
	deleteReq := httptest.NewRequest(http.MethodDelete, baseURL+fmt.Sprintf("api/v1/users/%d", userId), nil)
	deleteRecorder := httptest.NewRecorder()

	userImpl := &controller.UserImpl{DB: db}
	userImpl.DeleteApiV1UsersId(deleteRecorder, deleteReq, userId)

	actual := deleteRecorder.Code
	expected := http.StatusOK

	assert.Equal(t, expected, actual)
}

func TestPutApiV1UsersId(t *testing.T) {
	db := ContainerDB(t)
	_ = createUser(db, user1JSON)

	var userId int32 = 1
	putReq := httptest.NewRequest(http.MethodPut, baseURL+fmt.Sprintf("api/v1/users/%d", userId), strings.NewReader(user2JSON))
	putRecorder := httptest.NewRecorder()

	userImpl := &controller.UserImpl{DB: db}
	userImpl.PutApiV1UsersId(putRecorder, putReq, userId)

	var actual controller.User
	_ = json.NewDecoder(putRecorder.Body).Decode(&actual)

	assert.Equal(t, user2, actual)
}

func TestGetApiV1Users(t *testing.T) {
	db := ContainerDB(t)
	_ = createUser(db, user1JSON)
	_ = createUser(db, user2JSON)

	userImpl := &controller.UserImpl{DB: db}
	getReq := httptest.NewRequest(http.MethodGet, baseURL+"api/v1/users", nil)
	status := string(controller.Activated)
	params := controller.GetApiV1UsersParams{
		Current:  0,
		PageSize: 10,
		Status:   &status,
	}
	getRecorder := httptest.NewRecorder()
	userImpl.GetApiV1Users(getRecorder, getReq, params)

	var actual []controller.User
	_ = json.NewDecoder(getRecorder.Body).Decode(&actual)
	expected := []controller.User{
		user1,
		user2,
	}

	assert.Equal(t, expected, actual)
}

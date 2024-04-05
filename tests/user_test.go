package tests

import (
	"context"
	"encoding/json"
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

func TestPostApiV1Users(t *testing.T) {
	reqBodyJSON := `{
"username": "username1",
"password": "password1",
"email": "email@gamil.com1",
"phone": "18682635684",
"remark": "remark1",
"status": "activated",
"created": "2024-04-04 13:56:35.671521",
"updated": "2024-04-05 13:56:35.671521"
}`
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/users", strings.NewReader(reqBodyJSON))
	var reqBody controller.PostApiV1UsersJSONRequestBody
	_ = json.NewDecoder(strings.NewReader(reqBodyJSON)).Decode(&reqBody)
	w := httptest.NewRecorder()
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
	userImpl := &controller.UserImpl{DB: model.Setup(ctx, dsn)}
	userImpl.PostApiV1Users(w, req)
	var actual controller.User
	_ = json.NewDecoder(w.Body).Decode(&actual)
	name := ""
	email := "email@gamil.com1"
	phone := "18682635684"
	remark := "remark1"
	status := controller.Activated
	created := "2024-04-04 13:56:35.671521"
	updated := "2024-04-05 13:56:35.671521"
	expected := controller.User{
		Username: "username1",
		Name:     &name,
		Email:    &email,
		Phone:    &phone,
		Remark:   &remark,
		Status:   &status,
		Created:  &created,
		Updated:  &updated,
	}
	assert.Equal(t, expected, actual)
}

func TestGetApiV1UsersId(t *testing.T) {
	reqBodyJSON := `{
"username": "username1",
"password": "password1",
"name": "name1",
"email": "email@gamil.com1",
"phone": "18682635684",
"remark": "remark1",
"status": "activated",
"created": "2024-04-04 13:56:35.671521",
"updated": "2024-04-05 13:56:35.671521"
}`
	postReq := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/users", strings.NewReader(reqBodyJSON))
	var reqBody controller.PostApiV1UsersJSONRequestBody
	_ = json.NewDecoder(strings.NewReader(reqBodyJSON)).Decode(&reqBody)
	postRecorder := httptest.NewRecorder()
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
	userImpl := &controller.UserImpl{DB: model.Setup(ctx, dsn)}
	userImpl.PostApiV1Users(postRecorder, postReq)

	getReq := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/users/1", nil)
	var userId int32 = 1
	getRecorder := httptest.NewRecorder()
	userImpl.GetApiV1UsersId(getRecorder, getReq, userId)

	var actual controller.User
	_ = json.NewDecoder(getRecorder.Body).Decode(&actual)
	name := "name1"
	email := "email@gamil.com1"
	phone := "18682635684"
	remark := "remark1"
	status := controller.Activated
	created := "2024-04-04 13:56:35.671521"
	updated := "2024-04-05 13:56:35.671521"
	expected := controller.User{
		Username: "username1",
		Name:     &name,
		Email:    &email,
		Phone:    &phone,
		Remark:   &remark,
		Status:   &status,
		Created:  &created,
		Updated:  &updated,
	}
	assert.Equal(t, expected, actual)
}

func TestDeleteApiV1UsersId(t *testing.T) {
	reqBodyJSON := `{
"username": "username1",
"password": "password1",
"name": "name1",
"email": "email@gamil.com1",
"phone": "18682635684",
"remark": "remark1",
"status": "activated",
"created": "2024-04-04 13:56:35.671521",
"updated": "2024-04-05 13:56:35.671521"
}`
	postReq := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/users", strings.NewReader(reqBodyJSON))
	var reqBody controller.PostApiV1UsersJSONRequestBody
	_ = json.NewDecoder(strings.NewReader(reqBodyJSON)).Decode(&reqBody)
	postRecorder := httptest.NewRecorder()
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
	userImpl := &controller.UserImpl{DB: model.Setup(ctx, dsn)}
	userImpl.PostApiV1Users(postRecorder, postReq)

	deleteReq := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/api/v1/users/1", nil)
	deleteRecorder := httptest.NewRecorder()
	var userId int32 = 1
	userImpl.DeleteApiV1UsersId(deleteRecorder, deleteReq, userId)

	actual := deleteRecorder.Code
	expected := http.StatusOK
	assert.Equal(t, expected, actual)
}

func TestPutApiV1UsersId(t *testing.T) {
	reqBodyJSON := `{
"username": "username1",
"password": "password1",
"email": "email@gamil.com1",
"phone": "18682635684",
"remark": "remark1",
"status": "activated",
"created": "2024-04-04 13:56:35.671521",
"updated": "2024-04-05 13:56:35.671521"
}`
	postReq := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/users", strings.NewReader(reqBodyJSON))
	var reqBody controller.PostApiV1UsersJSONRequestBody
	_ = json.NewDecoder(strings.NewReader(reqBodyJSON)).Decode(&reqBody)
	postRecorder := httptest.NewRecorder()
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
	userImpl := &controller.UserImpl{DB: model.Setup(ctx, dsn)}
	userImpl.PostApiV1Users(postRecorder, postReq)

	putBodyJSON := `{
"username": "username2",
"password": "password2",
"email": "email@gamil.com2",
"phone": "2",
"remark": "remark2",
"status": "frozen",
"created": "2024-03-04 13:56:35.671521",
"updated": "2024-03-05 13:56:35.671521"
}`
	putReq := httptest.NewRequest(http.MethodPut, "http://localhost:8080/api/v1/users/1", strings.NewReader(putBodyJSON))
	putRecorder := httptest.NewRecorder()
	var userId int32 = 1
	userImpl.PutApiV1UsersId(putRecorder, putReq, userId)

	var actual controller.User
	_ = json.NewDecoder(putRecorder.Body).Decode(&actual)
	name := ""
	email := "email@gamil.com2"
	phone := "2"
	remark := "remark2"
	status := controller.Frozen
	created := "2024-03-04 13:56:35.671521"
	updated := "2024-03-05 13:56:35.671521"
	expected := controller.User{
		Username: "username2",
		Name:     &name,
		Email:    &email,
		Phone:    &phone,
		Remark:   &remark,
		Status:   &status,
		Created:  &created,
		Updated:  &updated,
	}
	assert.Equal(t, expected, actual)
}

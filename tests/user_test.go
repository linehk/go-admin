package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/linehk/go-admin/config"
	"github.com/linehk/go-admin/controller"
	"github.com/linehk/go-admin/model"
	"github.com/stretchr/testify/assert"
)

func TestPostApiV1Users(t *testing.T) {
	reqBody := "{\"password\":\"password\",\"username\":\"username\"}\n"
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/users", strings.NewReader(reqBody))
	w := httptest.NewRecorder()
	config.Setup()
	userImpl := &controller.UserImpl{DB: model.Setup()}
	userImpl.PostApiV1Users(w, req)
	actual := w.Body.String()
	expected := "{\"password\":\"password\",\"username\":\"username\"}\n"
	assert.Equal(t, expected, actual)
}

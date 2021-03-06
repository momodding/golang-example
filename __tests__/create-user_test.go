package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/saefullohmaslul/golang-example/src/app"
	"github.com/saefullohmaslul/golang-example/src/database"
	"github.com/saefullohmaslul/golang-example/src/global/types"
	"github.com/saefullohmaslul/golang-example/src/middlewares/exception"
	"github.com/stretchr/testify/assert"
)

func initTest(body map[string]interface{}) (*httptest.ResponseRecorder, *gin.Engine) {
	r := gin.Default()
	app := new(app.Application)
	app.CreateTest(r)

	w := httptest.NewRecorder()
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/user", strings.NewReader(string(b)))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	return w, r
}

func TestCreateUserSuccess(t *testing.T) {
	defer database.DropAllTable()
	body := map[string]interface{}{
		"name":     "User Test",
		"email":    "user@email.com",
		"password": "123456",
	}
	w, _ := initTest(body)

	actual := types.CreateUserResponse{}
	if err := json.Unmarshal([]byte(w.Body.String()), &actual); err != nil {
		panic(err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Success create user", actual.Message)
	assert.Equal(t, http.StatusOK, actual.Status)
	assert.NotEmpty(t, actual.Result)
}

func TestCreateUserExist(t *testing.T) {
	defer database.DropAllTable()
	body := map[string]interface{}{
		"name":     "User Test",
		"email":    "user@email.com",
		"password": "123456",
	}
	_, r := initTest(body)

	w := httptest.NewRecorder()
	b, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "/user", strings.NewReader(string(b)))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	actual := exception.Exception{}
	if err := json.Unmarshal([]byte(w.Body.String()), &actual); err != nil {
		panic(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, http.StatusBadRequest, actual.Status)
	assert.NotEmpty(t, actual.Errors)
	assert.Equal(t, "User already exist", actual.Errors.Message)
	assert.Equal(t, "USER_ALREADY_EXIST", actual.Errors.Flag)
}

func TestCreateUserInvalidBodyName(t *testing.T) {
	defer database.DropAllTable()
	body := map[string]interface{}{
		"name":     123,
		"email":    "user@email.com",
		"password": "123456",
	}
	w, _ := initTest(body)

	actual := exception.Exception{}
	if err := json.Unmarshal([]byte(w.Body.String()), &actual); err != nil {
		panic(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, http.StatusBadRequest, actual.Status)
	assert.NotEmpty(t, actual.Errors)
	assert.NotEmpty(t, actual.Errors.Message)
	assert.Equal(t, "INVALID_BODY", actual.Errors.Flag)
}

func TestCreateUserInvalidBodyEmail(t *testing.T) {
	defer database.DropAllTable()
	body := map[string]interface{}{
		"name":     "User Test",
		"email":    "useremail.com",
		"password": "123456",
	}
	w, _ := initTest(body)

	actual := exception.Exception{}
	if err := json.Unmarshal([]byte(w.Body.String()), &actual); err != nil {
		panic(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, http.StatusBadRequest, actual.Status)
	assert.NotEmpty(t, actual.Errors)
	assert.NotEmpty(t, actual.Errors.Message)
	assert.Equal(t, "INVALID_BODY", actual.Errors.Flag)
}

func TestCreateUserInvalidBodyPassword(t *testing.T) {
	defer database.DropAllTable()
	body := map[string]interface{}{
		"name":  "User Test",
		"email": "user@email.com",
	}
	w, _ := initTest(body)

	actual := exception.Exception{}
	if err := json.Unmarshal([]byte(w.Body.String()), &actual); err != nil {
		panic(err)
	}

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, http.StatusBadRequest, actual.Status)
	assert.NotEmpty(t, actual.Errors)
	assert.NotEmpty(t, actual.Errors.Message)
	assert.Equal(t, "INVALID_BODY", actual.Errors.Flag)
}

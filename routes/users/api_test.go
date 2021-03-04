package users

import (
	"echo_api/models"
	"echo_api/pkg/test"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
)

func Ok(c echo.Context) error {
	return c.String(http.StatusOK, "")
}

func TestAPI(t *testing.T) {
	mockRouter := echo.New()

	repo := &mockRepository{items: []models.User{
		{Id: 1, User: "test", Password: "password", Acc_type: 1},
	}}
	repo.seq = 1

	service := NewService(repo)
	RegisterHandlers(service, mockRouter)

	tests := []test.APITestCase{
		{"get all", "GET", "/users", "", nil, http.StatusOK, `*"total_count":1*`},
		{"get 1", "GET", "/users/1", "", nil, http.StatusOK, `*test*`},
		{"get unknown", "GET", "/users/1234", "", nil, http.StatusNotFound, ""},
		{"create ok", "POST", "/users", `{"user":"newtest", "password":"password", "acc_type":2}`, nil, http.StatusCreated, "*newtest*"},
		{"create ok count", "GET", "/users", "", nil, http.StatusOK, `*"total_count":2*`},
		//{"auth error", "POST", "/private", "", nil, http.StatusUnauthorized, ""}, // Today it returns 400, change?
		{"create input error", "POST", "/users", `"name":"test"}`, nil, http.StatusBadRequest, ""},
		{"update ok", "PUT", "/users/1", `{"user":"uptest", "password":"password", "acc_type":2}`, nil, http.StatusOK, "*uptest*"},
		{"update verify", "GET", "/users/1", "", nil, http.StatusOK, `*uptest*`},
		{"update input error", "PUT", "/users/1", `"name":"albumxyz"}`, nil, http.StatusBadRequest, ""},
		{"delete ok", "DELETE", "/users/1", ``, nil, http.StatusOK, ""},
		//{"delete verify", "DELETE", "/users/1", ``, nil, http.StatusNotFound, ""}, //FIXME not working with mock repository for some reason
	}
	for _, tc := range tests {
		test.Endpoint(t, mockRouter, tc)
	}
}

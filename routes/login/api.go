package login

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterHandlers(service Service, e *echo.Echo, m ...echo.MiddlewareFunc) {
	res := resource{service}
	e.POST("/login", res.login)
}

type resource struct {
	service Service
}

func (r resource) login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	t := r.service.Login(c, username, password)
	if len(t) == 0 {
		return c.String(http.StatusUnauthorized, "Username or Password incorrect\n")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

package users

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"echo_api/pkg/pagination"
)

//TODO: require jwt auth for some of the methods?
func RegisterHandlers(service Service, e *echo.Echo, m ...echo.MiddlewareFunc) {
	res := resource{service}

	e.POST("/users", res.create, m...)
	e.GET("/users", res.query, m...)
	e.GET("/users/:id", res.get, m...)
	e.PUT("/users/:id", res.update, m...)
	e.DELETE("/users/:id", res.delete, m...)
}

type resource struct {
	service Service
}

// func read_user(c echo.Context) error {
// 	id, _ := strconv.Atoi(c.Param("id"))

func (r resource) get(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := r.service.Get(c, id)
	if err == sql.ErrNoRows {
		return c.String(http.StatusNotFound, "")
	}
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func (r resource) query(c echo.Context) error {
	count, err := r.service.Count(c)
	if err != nil {
		return err
	}
	pages := pagination.NewFromRequest(c.Request(), count)
	users, err := r.service.Query(c, pages.Offset(), pages.Limit())
	if err != nil {
		return err
	}
	pages.Items = users
	return c.JSON(http.StatusOK, pages)
}

func (r resource) create(c echo.Context) error {
	user_request := &CreateUserRequest{}
	if err := c.Bind(user_request); err != nil {
		return c.String(http.StatusBadRequest, "Invalid JSON\n")
	}

	if len(user_request.User) == 0 || len(user_request.Password) == 0 {
		return c.String(http.StatusBadRequest, "Missing values in JSON\n")
	}

	if user_request.Acc_type < 0 || user_request.Acc_type > 4 {
		return c.String(http.StatusBadRequest, "Invalid Value for acc_type\n")
	}

	new_user, err := r.service.Create(c, *user_request)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, new_user)
}

func (r resource) update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Id is an int\n")
	}
	user_request := &UpdateUserRequest{}
	if err := c.Bind(user_request); err != nil {
		return c.String(http.StatusBadRequest, "Invalid JSON\n")
	}

	if len(user_request.User) == 0 || len(user_request.Password) == 0 {
		return c.String(http.StatusBadRequest, "Missing values in JSON\n")
	}

	if user_request.Acc_type < 0 || user_request.Acc_type > 4 {
		return c.String(http.StatusBadRequest, "Invalid Value for acc_type\n")
	}

	new_user, err := r.service.Update(c, id, *user_request)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, new_user)
}

func (r resource) delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Id is an int\n")
	}

	err = r.service.Delete(c, id)
	if err == sql.ErrNoRows {
		return c.String(http.StatusNotFound, "")
	}
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "")

}

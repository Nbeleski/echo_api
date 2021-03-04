package maps

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterHandlers(service Service, e *echo.Echo, m ...echo.MiddlewareFunc) {
	res := resource{service}
	e.POST("/maps/directions", res.directions)
}

type resource struct {
	service Service
}

func (r resource) directions(c echo.Context) error {
	dirReq := &DirectionsRequest{}
	if err := c.Bind(dirReq); err != nil {
		return c.String(http.StatusBadRequest, "Invalid JSON\n")
	}
	if len(dirReq.Origin) == 0 || len(dirReq.Destination) == 0 {
		return c.String(http.StatusBadRequest, "Missing values in JSON\n")
	}

	route, err := r.service.Directions(c, *dirReq)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, route)
}

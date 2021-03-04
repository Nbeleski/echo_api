package maps

import (
	"context"

	"github.com/labstack/echo/v4"
	"googlemaps.github.io/maps"
)

type Service interface {
	Directions(ctx echo.Context, req DirectionsRequest) ([]maps.Route, error)
}

type service struct {
	APIKey string
}

func NewService(APIKey string) Service {
	return service{APIKey}
}

type DirectionsRequest struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
}

func (s service) Directions(ctx echo.Context, req DirectionsRequest) ([]maps.Route, error) {
	c, err := maps.NewClient(maps.WithAPIKey(s.APIKey))
	if err != nil {
		return []maps.Route{}, err
	}
	r := &maps.DirectionsRequest{
		Origin:      req.Origin,
		Destination: req.Destination,
	}
	route, _, err := c.Directions(context.Background(), r)
	if err != nil {
		return []maps.Route{}, err
	}

	return route, nil
}

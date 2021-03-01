package auth

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
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

func GenareteSaltedPassword(password []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12) //minCost = 12, safer than the default 4
	return hash, err
}

func ComparePasswords(hashedPwd []byte, plainPwd []byte) bool { // Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	err := bcrypt.CompareHashAndPassword(hashedPwd, plainPwd)
	if err != nil {
		//log.Println(err)
		return false
	}
	return true
}

var IsLoggedIn = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: []byte(secret),
})

func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		acc_type := claims["account_type"].(float64)
		if acc_type != 4 {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}

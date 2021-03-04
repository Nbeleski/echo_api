package authentication

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	IsLoggedIn() echo.MiddlewareFunc
	IsAdmin(next echo.HandlerFunc) echo.HandlerFunc
}

type service struct {
	secret string
}

func NewService(secret string) Service {
	return service{secret}
}

// var IsLoggedIn = middleware.JWTWithConfig(middleware.JWTConfig{
// 	SigningKey: []byte(secret),
// })

func (s service) IsLoggedIn() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(s.secret),
	})
}

func (s service) IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
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

func GenareteSaltedPassword(password []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12) //minCost = 12, safer than the default 4
	return hash, err
}

func ComparePasswords(hashedPwd []byte, plainPwd []byte) bool { // Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	err := bcrypt.CompareHashAndPassword(hashedPwd, plainPwd)
	if err != nil {
		return false
	}
	return true
}

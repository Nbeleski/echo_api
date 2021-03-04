package login

import (
	"echo_api/pkg/authentication"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type Service interface {
	Login(ctx echo.Context, user, password string) string
}

type service struct {
	repo   Repository
	secret string
}

func NewService(repo Repository, secret string) Service {
	return service{repo, secret}
}

func (s service) Login(ctx echo.Context, form_user, form_password string) string {
	user, err := s.repo.GetByUser(ctx, form_user)
	t := ""
	if err != nil {
		return t
	}

	if authentication.ComparePasswords([]byte(user.Password), []byte(form_password)) {
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = user.User
		claims["account_type"] = user.Acc_type
		claims["exp"] = time.Now().Add(time.Hour * 720).Unix()

		// Generate encoded token and send it as response.
		t, err = token.SignedString([]byte(s.secret))
		if err != nil {
			return ""
		}
	}

	return t
}

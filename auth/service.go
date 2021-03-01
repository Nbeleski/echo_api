package auth

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

var secret = "OURAPI-f0ab801e-750d-11eb-9439-0242ac130002"

type Service interface {
	Login(ctx echo.Context, user, password string) string
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return service{repo}
}

func (s service) Login(ctx echo.Context, form_user, form_password string) string {
	user, err := s.repo.GetByUser(ctx, form_user)
	t := ""
	if err != nil {
		return t
	}

	if ComparePasswords([]byte(user.Password), []byte(form_password)) {
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = user.User
		claims["account_type"] = user.Acc_type
		claims["exp"] = time.Now().Add(time.Hour * 720).Unix()

		// Generate encoded token and send it as response.
		t, err = token.SignedString([]byte(secret))
		if err != nil {
			return ""
		}
	}

	return t
}

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// retrieve password for username
	database, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Failed to open connection to Db")
	}
	statement, err := database.Prepare("SELECT password, acc_type FROM tab_users WHERE user=?")
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "")
	}
	rows, err := statement.Query(username)

	if err != nil {
		return c.String(http.StatusInternalServerError, "Query Failed")
	}

	var count = 0
	var saved_password string
	var acc_type int
	for rows.Next() {
		rows.Scan(&saved_password, &acc_type)
		count += 1
	}

	if count == 0 {
		return c.String(http.StatusUnauthorized, "Username or Password incorrect")
	}

	if !ComparePasswords([]byte(saved_password), []byte(password)) {
		return c.String(http.StatusUnauthorized, "Username or Password incorrect")
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = username
	claims["account_type"] = acc_type
	claims["exp"] = time.Now().Add(time.Hour * 720).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

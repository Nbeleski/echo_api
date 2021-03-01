package main

import (
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"echo_api/auth"
	"echo_api/users"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339_nano}][${remote_ip}][${host}]:uri=${uri}, method=${method}, status=${status}\n", //, error:${error}\n",
	}))

	// Middleware to debug POST requests:
	// e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
	// 	log.Println(string(reqBody))
	// }))

	login_service := auth.NewService(auth.NewRepository("sqlite3", "test.db"))
	users_service := users.NewService(users.NewRepository("sqlite3", "test.db"))

	// Routes
	auth.RegisterHandlers(login_service, e)
	users.RegisterHandlers(users_service, e)

	// Example of private and administrator routes
	e.GET("/private", private, auth.IsLoggedIn)
	e.GET("/admin", private, auth.IsLoggedIn, auth.IsAdmin)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func private(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	acc_type := claims["account_type"].(float64)
	if acc_type == 4 {
		return c.String(http.StatusOK, "Welcome Admin "+name+"!"+" Your permission level is "+strconv.FormatFloat(acc_type, 'g', 1, 64))
	}
	return c.String(http.StatusOK, "Welcome "+name+"!"+" Your permission level is "+strconv.FormatFloat(acc_type, 'g', 1, 64))
}

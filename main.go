package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"echo_api/pkg/authentication"
	"echo_api/pkg/config"
	"echo_api/routes/login"
	"echo_api/routes/maps"
	"echo_api/routes/users"
)

func main() {
	config, err := config.LoadConfiguration("config.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	db, err := sql.Open(config.Database.Driver, config.Database.Datasource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Echo instance
	e := echo.New()

	// Static files
	e.Static("/static", "assets")

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339_nano}][${remote_ip}][${host}]:uri=${uri}, method=${method}, status=${status}\n", //, error:${error}\n",
	}))

	// Middleware to debug POST requests:
	// e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
	// 	log.Println(string(reqBody))
	// }))

	auth := authentication.NewService(config.JWT_secret)
	login_service := login.NewService(login.NewRepository(db), config.JWT_secret)
	users_service := users.NewService(users.NewRepository(db))
	maps_service := maps.NewService(config.MapsAPIKey)

	// Routes
	login.RegisterHandlers(login_service, e)
	users.RegisterHandlers(users_service, e, auth.IsLoggedIn()) //login only required for POST/PUT/DELETE
	maps.RegisterHandlers(maps_service, e)

	// Example of private and administrator routes
	e.GET("/private", private, auth.IsLoggedIn())
	e.GET("/admin", private, auth.IsLoggedIn(), auth.IsAdmin)

	e.File("/", "static/index.html")

	// Start server
	e.Logger.Fatal(e.Start(config.Host + ":" + config.Port))
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

package auth

import (
	"database/sql"
	"echo_api/models"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	GetByUser(ctx echo.Context, user string) (models.User, error)
}

type repository struct {
	dbtype string
	dbfile string
}

// NewRepository creates a new album repository
func NewRepository(dbtype, dbfile string) Repository {
	return repository{dbtype, dbfile}
}

func (r repository) GetByUser(c echo.Context, user string) (models.User, error) {
	var saved_user models.User
	database, err := sql.Open(r.dbtype, r.dbfile)
	if err != nil {
		return saved_user, err
	}
	statement, err := database.Prepare("SELECT password, acc_type FROM tab_users WHERE user=?")
	if err != nil {
		return saved_user, err
	}
	rows, err := statement.Query(user)
	if err != nil {
		return saved_user, err
	}

	for rows.Next() {
		rows.Scan(&saved_user.Password, &saved_user.Acc_type)
	}

	return saved_user, nil
}

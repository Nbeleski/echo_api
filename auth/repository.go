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
	db *sql.DB
}

// NewRepository creates a new album repository
func NewRepository(db *sql.DB) Repository {
	return repository{db}
}

func (r repository) GetByUser(c echo.Context, user string) (models.User, error) {
	var saved_user models.User
	statement, err := r.db.Prepare("SELECT password, acc_type FROM tab_users WHERE user=?")
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

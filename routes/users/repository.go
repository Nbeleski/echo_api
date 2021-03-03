package users

import (
	"database/sql"
	"echo_api/models"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Get(ctx echo.Context, id int) (models.User, error)
	Count(ctx echo.Context) (int, error)
	Query(ctx echo.Context, offset, limit int) ([]models.User, error)
	Create(ctx echo.Context, user models.User) (int, error)
	Update(ctx echo.Context, user models.User) error
	Delete(ctx echo.Context, id int) error
}

type repository struct {
	db *sql.DB
}

// NewRepository creates a new album repository
func NewRepository(db *sql.DB) Repository {
	return repository{db}
}

func (r repository) Get(c echo.Context, id int) (models.User, error) {
	var user models.User
	statement, err := r.db.Prepare("SELECT user, password, acc_type FROM tab_users WHERE id=?")
	if err != nil {
		return user, err
	}
	rows, err := statement.Query(id)
	if err != nil {
		return user, err
	}

	count := 0
	for rows.Next() {
		user.Id = id
		rows.Scan(&user.User, &user.Password, &user.Acc_type)
		count++
	}
	if count == 0 {
		return user, sql.ErrNoRows
	}

	return user, err
}

func (r repository) Query(c echo.Context, offset, limit int) ([]models.User, error) {
	var users []models.User
	statement, err := r.db.Prepare("SELECT t.id, t.user, t.password, t.acc_type FROM (SELECT * FROM tab_users LIMIT ? OFFSET ?) t ORDER BY id COLLATE NOCASE")
	if err != nil {
		return users, err
	}
	rows, err := statement.Query(limit, offset)
	if err != nil {
		return users, err
	}

	var id, acc_type int
	var user, password string
	for rows.Next() {
		rows.Scan(&id, &user, &password, &acc_type)
		users = append(users, models.User{
			Id:       id,
			User:     user,
			Password: password,
			Acc_type: acc_type,
		})
	}

	return users, nil
}

func (r repository) Create(c echo.Context, user models.User) (int, error) {
	statement, err := r.db.Prepare("INSERT INTO tab_users VALUES (null, ?, ? ,?)")
	if err != nil {
		return -1, err
	}
	result, err := statement.Exec(user.User, user.Password, user.Acc_type)
	if err != nil {
		return -1, err
	}

	new_id, _ := result.LastInsertId()
	return int(new_id), nil
}

func (r repository) Update(c echo.Context, user models.User) error {
	statement, err := r.db.Prepare("UPDATE tab_users SET user=?, password=?, acc_type=? WHERE id=?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(user.User, user.Password, user.Acc_type, user.Id)
	return err
}

func (r repository) Delete(c echo.Context, id int) error {
	statement, err := r.db.Prepare("DELETE FROM tab_users WHERE id=?")
	if err != nil {
		return err
	}
	res, err := statement.Exec(id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if count > 0 {
		return nil
	}
	return sql.ErrNoRows
}

func (r repository) Count(c echo.Context) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM tab_users").Scan(&count)
	return count, err
}

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
	Delete(ctx echo.Context, id int) (bool, error)
}

type repository struct {
	dbtype string
	dbfile string
}

// NewRepository creates a new album repository
func NewRepository(dbtype, dbfile string) Repository {
	return repository{dbtype, dbfile}
}

func (r repository) Get(c echo.Context, id int) (models.User, error) {
	var user models.User
	database, err := sql.Open(r.dbtype, r.dbfile)
	if err != nil {
		return user, err
	}
	statement, err := database.Prepare("SELECT user, password, acc_type FROM tab_users WHERE id=?")
	if err != nil {
		return user, err
	}
	rows, err := statement.Query(id)
	if err != nil {
		return user, err
	}

	for rows.Next() {
		user.Id = id
		rows.Scan(&user.User, &user.Password, &user.Acc_type)
	}

	return user, err
}

func (r repository) Query(c echo.Context, offset, limit int) ([]models.User, error) {
	var users []models.User
	database, err := sql.Open(r.dbtype, r.dbfile)
	if err != nil {
		return users, err
	}
	statement, err := database.Prepare("SELECT t.id, t.user, t.password, t.acc_type FROM (SELECT * FROM tab_users LIMIT ? OFFSET ?) t ORDER BY id COLLATE NOCASE")
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
	database, err := sql.Open(r.dbtype, r.dbfile)
	if err != nil {
		return -1, err
	}
	statement, err := database.Prepare("INSERT INTO tab_users VALUES (null, ?, ? ,?)")
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
	database, err := sql.Open(r.dbtype, r.dbfile)
	if err != nil {
		return err
	}
	statement, err := database.Prepare("UPDATE tab_users SET user=?, password=?, acc_type=? WHERE id=?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(user.User, user.Password, user.Acc_type, user.Id)
	return err
}

func (r repository) Delete(c echo.Context, id int) (bool, error) {
	database, err := sql.Open(r.dbtype, r.dbfile)
	if err != nil {
		return false, err
	}
	statement, err := database.Prepare("DELETE FROM tab_users WHERE id=?")
	if err != nil {
		return false, err
	}
	res, err := statement.Exec(id)
	if err != nil {
		return false, err
	}

	count, err := res.RowsAffected()
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (r repository) Count(c echo.Context) (int, error) {
	database, err := sql.Open(r.dbtype, r.dbfile)
	if err != nil {
		return -1, err
	}
	var count int
	err = database.QueryRow("SELECT COUNT(*) FROM tab_users").Scan(&count)
	return count, err
}

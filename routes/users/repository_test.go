package users

import (
	"database/sql"
	"echo_api/models"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	// create temporary test table
	os.Create("users_repo_test.db")
	db, err := sql.Open("sqlite3", "users_repo_test.db")
	assert.Nil(t, err)
	_, err = db.Exec("CREATE TABLE tab_users (id INTEGER PRIMARY KEY UNIQUE, user TEXT NOT NULL UNIQUE, password TEXT NOT NULL, acc_type INTEGER DEFAULT (0) CHECK (acc_type >= 0 AND acc_type <= 4));")
	assert.Nil(t, err)

	repo := NewRepository(db)
	var ctx TestContext

	// initial count
	count, err := repo.Count(ctx)
	assert.Nil(t, err)

	// create
	id, err := repo.Create(ctx, models.User{
		Id:       0,
		User:     "user1",
		Password: "password1",
		Acc_type: 1,
	})
	assert.Equal(t, 1, id)
	assert.Nil(t, err)
	count2, _ := repo.Count(ctx)
	assert.Equal(t, 1, count2-count)

	// get
	user, err := repo.Get(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "user1", user.User)
	assert.Equal(t, 1, user.Acc_type)
	_, err = repo.Get(ctx, 2)
	assert.Equal(t, sql.ErrNoRows, err)

	// update
	err = repo.Update(ctx, models.User{
		Id:       1,
		User:     "updated user1",
		Password: "password1",
		Acc_type: 1,
	})
	assert.Nil(t, err)
	user, err = repo.Get(ctx, id)
	assert.Equal(t, "updated user1", user.User)

	// query
	users_list, err := repo.Query(ctx, 0, count2)
	assert.Nil(t, err)
	assert.Equal(t, count2, len(users_list))

	// delete
	err = repo.Delete(ctx, 1)
	assert.Nil(t, err)
	_, err = repo.Get(ctx, 1)
	assert.Equal(t, sql.ErrNoRows, err)
	err = repo.Delete(ctx, 1)
	assert.Equal(t, sql.ErrNoRows, err)

	// remove temporary test table
	os.Remove("users_repo_test.db")
}

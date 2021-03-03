package users

import (
	"database/sql"
	"echo_api/models"
	"errors"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_service_CRUD(t *testing.T) {
	mockRepo := &mockRepository{}
	mockRepo.seq = 1
	s := NewService(mockRepo)

	var ctx TestContext

	// initial count
	count, _ := s.Count(ctx)
	assert.Equal(t, 0, count)

	// create
	user, err := s.Create(ctx, CreateUserRequest{
		User:     "test",
		Password: "1234",
		Acc_type: 1,
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, user.Id)
	id := user.Id
	assert.Equal(t, id, user.Id)
	assert.Equal(t, "test", user.User)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// validation length error in create
	_, err = s.Create(ctx, CreateUserRequest{User: "t", Password: "1234", Acc_type: 1})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// validation required error in create
	_, err = s.Create(ctx, CreateUserRequest{User: "valid"})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	user, _ = s.Create(ctx, CreateUserRequest{
		User:     "test2",
		Password: "1234",
		Acc_type: 1,
	})
	id = user.Id

	// update
	user, err = s.Update(ctx, id, UpdateUserRequest{
		User:     "test updated",
		Password: "1234",
		Acc_type: 1,
	})
	assert.Nil(t, err)
	assert.Equal(t, "test updated", user.User)

	// validation length error in update
	_, err = s.Update(ctx, id, UpdateUserRequest{User: "t", Password: "1234", Acc_type: 1})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 2, count)

	// validation required error in update
	_, err = s.Update(ctx, id, UpdateUserRequest{User: "valid"})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 2, count)

	// get
	user, err = s.Get(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "test updated", user.User)

	// query
	albums, _ := s.Query(ctx, 0, 0)
	assert.Equal(t, 2, len(albums))

	// delete
	err = s.Delete(ctx, id)
	assert.Nil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)
	err = s.Delete(ctx, id)
	assert.Equal(t, sql.ErrNoRows, err)
}

type mockRepository struct {
	items []models.User
	seq   int
}

type TestContext struct {
	echo.Context
}

var errCRUD = errors.New("error crud")

func (m mockRepository) Get(ctx echo.Context, id int) (models.User, error) {
	for _, item := range m.items {
		if item.Id == id {
			return item, nil
		}
	}
	return models.User{}, sql.ErrNoRows
}

func (m mockRepository) Count(ctx echo.Context) (int, error) {
	return len(m.items), nil
}

func (m mockRepository) Query(ctx echo.Context, offset, limit int) ([]models.User, error) {
	return m.items, nil
}

func (m *mockRepository) Create(ctx echo.Context, user models.User) (int, error) {
	if user.User == "error" {
		return 0, errCRUD
	}
	user.Id = m.seq
	m.seq++
	m.items = append(m.items, user)
	return len(m.items), nil
}

func (m *mockRepository) Update(ctx echo.Context, User models.User) error {
	if User.User == "error" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.Id == User.Id {
			m.items[i] = User
			break
		}
	}
	return nil
}

func (m *mockRepository) Delete(ctx echo.Context, id int) error {
	for i, item := range m.items {
		if item.Id == id {
			m.items[i] = m.items[len(m.items)-1]
			m.items = m.items[:len(m.items)-1]
			m.seq--
			return nil
		}
	}
	return sql.ErrNoRows
}

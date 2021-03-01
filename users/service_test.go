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
	s := NewService(&mockRepository{})

	var ctx TestContext

	// initial count
	count, _ := s.Count(ctx)
	assert.Equal(t, 0, count)
}

type mockRepository struct {
	items []models.User
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

func (m *mockRepository) Create(ctx echo.Context, User models.User) (int, error) {
	if User.User == "error" {
		return 0, errCRUD
	}
	m.items = append(m.items, User)
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

func (m *mockRepository) Delete(ctx echo.Context, id int) (bool, error) {
	for i, item := range m.items {
		if item.Id == id {
			m.items[i] = m.items[len(m.items)-1]
			m.items = m.items[:len(m.items)-1]
			return true, nil
		}
	}
	return false, nil
}

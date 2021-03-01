package users

import (
	"echo_api/auth"
	"echo_api/models"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

type CreateUserRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Acc_type int    `json:"acc_type"`
}

type UpdateUserRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Acc_type int    `json:"acc_type"`
}

type Service interface {
	Get(ctx echo.Context, id int) (models.User, error)
	Query(ctx echo.Context, offset, limit int) ([]models.User, error)
	Create(ctx echo.Context, req CreateUserRequest) (models.User, error)
	Update(ctx echo.Context, id int, req UpdateUserRequest) (models.User, error)
	Delete(ctx echo.Context, id int) (bool, error)
	Count(ctx echo.Context) (int, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return service{repo}
}

// Get returns the album with the specified the album ID.
func (s service) Get(ctx echo.Context, id int) (models.User, error) {
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		return models.User{}, err
	}
	return user, nil //return models.User{user}?
}

func (s service) Query(ctx echo.Context, offset, limit int) ([]models.User, error) {
	return s.repo.Query(ctx, offset, limit)
}

func (s service) Create(ctx echo.Context, req CreateUserRequest) (models.User, error) {
	hash, _ := auth.GenareteSaltedPassword([]byte(req.Password))
	new_user := models.User{
		Id:       0,
		User:     req.User,
		Password: string(hash), //user_request.Password,
		Acc_type: req.Acc_type,
	}
	new_id, err := s.repo.Create(ctx, new_user)
	if err != nil {
		return models.User{}, err
	}
	new_user.Id = new_id
	return new_user, nil
}

func (s service) Update(ctx echo.Context, id int, req UpdateUserRequest) (models.User, error) {
	hash, _ := auth.GenareteSaltedPassword([]byte(req.Password))
	new_user := models.User{
		Id:       id,
		User:     req.User,
		Password: string(hash), //user_request.Password,
		Acc_type: req.Acc_type,
	}
	err := s.repo.Update(ctx, new_user)
	if err != nil {
		return models.User{}, err
	}

	return new_user, nil
}

func (s service) Delete(ctx echo.Context, id int) (bool, error) {
	return s.repo.Delete(ctx, id)
}

func (s service) Count(c echo.Context) (int, error) {
	return s.repo.Count(c)
}

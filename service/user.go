package service

import (
	"context"
	"errors"
	"github.com/cooljeffrey/petstore/model"
	"github.com/go-kit/kit/log"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) error
	CreateUsersWithArray(ctx context.Context, array []*model.User) error
	CreateUsersWithList(ctx context.Context, list []*model.User) error
	Login(ctx context.Context, username, password string) error
	Logout(ctx context.Context) error
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	UpdateUserByUsername(ctx context.Context, username string, user *model.User) error
	DeleteUserByUsername(ctx context.Context, username string) error
}

type userService struct {
	logger  log.Logger
	storage model.Storage
}

func NewUserService(logger log.Logger, storage model.Storage) UserService {
	return &userService{
		logger:  logger,
		storage: storage,
	}
}

func (s userService) CreateUser(ctx context.Context, user *model.User) error {
	return s.storage.CreateUser(user)
}

func (s userService) CreateUsersWithArray(ctx context.Context, array []*model.User) error {
	return s.storage.CreateManyUsers(array)
}

func (s userService) CreateUsersWithList(ctx context.Context, list []*model.User) error {
	return s.storage.CreateManyUsers(list)
}

func (s userService) Login(ctx context.Context, username, password string) error {
	user, err := s.storage.RetrieveUserByUsername(username)
	if err != nil {
		return err
	}
	if user.Password == password {
		// TODO implement token authentiation and authorization ?
		return nil
	} else {
		return errors.New("invalid username and password combination")
	}
}

func (s userService) Logout(ctx context.Context) error {
	// TODO implement token authentiation and authorization ?
	return nil
}

func (s userService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return s.storage.RetrieveUserByUsername(username)
}

func (s userService) UpdateUserByUsername(ctx context.Context, username string, user *model.User) error {
	_, err := s.storage.UpdateUserByUsername(username, user)
	if err != nil {
		return err
	}
	return nil
}

func (s userService) DeleteUserByUsername(ctx context.Context, username string) error {
	return s.storage.DeleteUserByUsername(username)
}

// Package service implements the business logic for user-service.
package service

import (
	"context"
	"errors"

	"furab-backend/services/user-service/internal/model"
	"furab-backend/services/user-service/internal/repository"
)

type User = model.User

type CreateUserRequest struct {
	UserID string
	Name   string
	Email  string
	Phone  string
}

type CreateUserResponse struct {
	UserID  string
	Message string
}

type UpdateUserRequest struct {
	Name  string
	Email string
}

// UserService defines the interface for user-service business logic.
type UserService interface {
	GetProfile(ctx context.Context) error
	UpdateProfile(ctx context.Context) error
	AddAddress(ctx context.Context) error
	DeleteAddress(ctx context.Context) error

	CreateUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error)
	GetUser(ctx context.Context, userID string) (*User, error)
	UpdateUser(ctx context.Context, userID string, req UpdateUserRequest) error
	DeactivateUser(ctx context.Context, userID string) error
}

// userServiceImpl is the concrete implementation of UserService.
type userServiceImpl struct {
	repo repository.UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(repo repository.UserRepository) UserService {
	return &userServiceImpl{repo: repo}
}

func (s *userServiceImpl) GetProfile(ctx context.Context) error { return nil }
func (s *userServiceImpl) UpdateProfile(ctx context.Context) error { return nil }
func (s *userServiceImpl) AddAddress(ctx context.Context) error { return nil }
func (s *userServiceImpl) DeleteAddress(ctx context.Context) error { return nil }

func (s *userServiceImpl) CreateUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error) {
	if req.Email == "" {
		return nil, errors.New("email required")
	}
	if req.Name == "" {
		return nil, errors.New("validation error")
	}

	user := &User{
		UserID: req.UserID,
		Name:   req.Name,
		Email:  req.Email,
		Phone:  req.Phone,
	}

	if err := s.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	return &CreateUserResponse{UserID: user.UserID, Message: "sukses"}, nil
}

func (s *userServiceImpl) GetUser(ctx context.Context, userID string) (*User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *userServiceImpl) UpdateUser(ctx context.Context, userID string, req UpdateUserRequest) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	user.Name = req.Name
	user.Email = req.Email
	return s.repo.Update(ctx, user)
}

func (s *userServiceImpl) DeactivateUser(ctx context.Context, userID string) error {
	_, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	return s.repo.Deactivate(ctx, userID)
}

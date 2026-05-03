package service

import (
	"context"
	"errors"

	"furab-backend/services/user-service/internal/model"
	"furab-backend/services/user-service/internal/repository"
)

type User = model.User

// ✅ Define error biar konsisten & profesional
var (
	ErrValidation   = errors.New("validation error")
	ErrUserNotFound = errors.New("user not found")
)

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

type userServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userServiceImpl{repo: repo}
}

// =======================
// UNUSED (boleh kasih TODO)
// =======================
func (s *userServiceImpl) GetProfile(ctx context.Context) error {
	// TODO: implement
	return nil
}
func (s *userServiceImpl) UpdateProfile(ctx context.Context) error {
	// TODO: implement
	return nil
}
func (s *userServiceImpl) AddAddress(ctx context.Context) error {
	// TODO: implement
	return nil
}
func (s *userServiceImpl) DeleteAddress(ctx context.Context) error {
	// TODO: implement
	return nil
}

// =======================
// CORE LOGIC
// =======================

func (s *userServiceImpl) CreateUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error) {
	// ✅ Validasi dirapikan (konsisten)
	if req.Email == "" || req.Name == "" {
		return nil, ErrValidation
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

	return &CreateUserResponse{
		UserID:  user.UserID,
		Message: "sukses",
	}, nil
}

func (s *userServiceImpl) GetUser(ctx context.Context, userID string) (*User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *userServiceImpl) UpdateUser(ctx context.Context, userID string, req UpdateUserRequest) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	user.Name = req.Name
	user.Email = req.Email

	return s.repo.Update(ctx, user)
}

func (s *userServiceImpl) DeactivateUser(ctx context.Context, userID string) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	return s.repo.Deactivate(ctx, userID)
}
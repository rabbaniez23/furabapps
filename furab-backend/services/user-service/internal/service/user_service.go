package service

import (
	"context"
	"errors"
	"strings"

	"furab-backend/services/user-service/internal/model"
	"furab-backend/services/user-service/internal/repository"
)

type User = model.User

// ✅ Define error biar konsisten & profesional
var (
	ErrValidation     = errors.New("validation error")
	ErrUserIDRequired = errors.New("user id required")
	ErrNameRequired   = errors.New("name required")
	ErrEmailRequired  = errors.New("email required")
	ErrPhoneRequired  = errors.New("phone required")
	ErrUserNotFound   = errors.New("user not found")
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

func normalizeInput(v string) string {
	return strings.TrimSpace(v)
}

func validateUserID(userID string) error {
	if normalizeInput(userID) == "" {
		return ErrUserIDRequired
	}
	return nil
}

func validateCreateUserRequest(req CreateUserRequest) error {
	if err := validateUserID(req.UserID); err != nil {
		return err
	}
	if normalizeInput(req.Name) == "" {
		return ErrNameRequired
	}
	if normalizeInput(req.Email) == "" {
		return ErrEmailRequired
	}
	if normalizeInput(req.Phone) == "" {
		return ErrPhoneRequired
	}
	return nil
}

func validateUpdateUserRequest(req UpdateUserRequest) error {
	if normalizeInput(req.Name) == "" {
		return ErrNameRequired
	}
	if normalizeInput(req.Email) == "" {
		return ErrEmailRequired
	}
	return nil
}

// =======================
// UNUSED (boleh kasih TODO)
// =======================
func (s *userServiceImpl) GetProfile(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return errors.New("not implemented")
}
func (s *userServiceImpl) UpdateProfile(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return errors.New("not implemented")
}
func (s *userServiceImpl) AddAddress(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return errors.New("not implemented")
}
func (s *userServiceImpl) DeleteAddress(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return errors.New("not implemented")
}

// =======================
// CORE LOGIC
// =======================

func (s *userServiceImpl) CreateUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	req.UserID = normalizeInput(req.UserID)
	req.Name = normalizeInput(req.Name)
	req.Email = normalizeInput(req.Email)
	req.Phone = normalizeInput(req.Phone)
	if err := validateCreateUserRequest(req); err != nil {
		return nil, err
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
		Message: "user created",
	}, nil
}

func (s *userServiceImpl) GetUser(ctx context.Context, userID string) (*User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	userID = normalizeInput(userID)
	if err := validateUserID(userID); err != nil {
		return nil, err
	}

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
	if err := ctx.Err(); err != nil {
		return err
	}
	userID = normalizeInput(userID)
	if err := validateUserID(userID); err != nil {
		return err
	}
	req.Name = normalizeInput(req.Name)
	req.Email = normalizeInput(req.Email)
	if err := validateUpdateUserRequest(req); err != nil {
		return err
	}

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
	if err := ctx.Err(); err != nil {
		return err
	}
	userID = normalizeInput(userID)
	if err := validateUserID(userID); err != nil {
		return err
	}

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	return s.repo.Deactivate(ctx, userID)
}
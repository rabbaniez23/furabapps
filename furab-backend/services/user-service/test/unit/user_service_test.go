package unit

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	// TODO: Sesuaikan import path ini dengan struktur project Anda setelah interface & model diperbarui
	// "furab-backend/services/user-service/internal/model"
	// "furab-backend/services/user-service/internal/service"
	// mock_repository "furab-backend/services/user-service/internal/mock" 
)

// ============================================================================
// Catatan: Struct dan interface di bawah ini adalah representasi dari spesifikasi.
// Dalam implementasi nyata, ini berada di package model/service dan mock digenerate oleh mockgen.
// ============================================================================

type User struct {
	UserID string
	Name   string
	Email  string
	Phone  string
	Status string
}

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

// Representasi dari UserService yang akan di-test
type UserService interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error)
	GetUser(ctx context.Context, userID string) (*User, error)
	UpdateUser(ctx context.Context, userID string, req UpdateUserRequest) error
	DeactivateUser(ctx context.Context, userID string) error
}

// Representasi manual dari MockUserRepository (seharusnya digenerate otomatis oleh gomock)
// Contoh command: mockgen -source=internal/repository/user_repository.go -destination=internal/mock/user_repository_mock.go -package=mock_repository
type MockUserRepository struct {
	SaveFunc       func(ctx context.Context, user *User) error
	FindByIDFunc   func(ctx context.Context, userID string) (*User, error)
	UpdateFunc     func(ctx context.Context, user *User) error
	DeactivateFunc func(ctx context.Context, userID string) error
}

func (m *MockUserRepository) Save(ctx context.Context, user *User) error { return m.SaveFunc(ctx, user) }
func (m *MockUserRepository) FindByID(ctx context.Context, userID string) (*User, error) { return m.FindByIDFunc(ctx, userID) }
func (m *MockUserRepository) Update(ctx context.Context, user *User) error { return m.UpdateFunc(ctx, user) }
func (m *MockUserRepository) Deactivate(ctx context.Context, userID string) error { return m.DeactivateFunc(ctx, userID) }

// Implementasi Dummy Service untuk Testing (Seharusnya dari internal/service)
type userServiceImpl struct {
	repo *MockUserRepository
}

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


// ============================================================================
// UNIT TESTS MULAI DARI SINI
// ============================================================================

func TestUserService_CreateUser(t *testing.T) {
	// gomock controller digunakan untuk mock standard, di sini kita pakai dummy struct untuk demo
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success - User berhasil dibuat", func(t *testing.T) {
		repo := &MockUserRepository{
			SaveFunc: func(ctx context.Context, user *User) error {
				if user.UserID != "1" {
					t.Errorf("Expected UserID 1, got %s", user.UserID)
				}
				return nil
			},
		}
		service := &userServiceImpl{repo: repo}

		req := CreateUserRequest{
			UserID: "1",
			Name:   "Erv",
			Email:  "erv@mail.com",
			Phone:  "08123",
		}

		res, err := service.CreateUser(context.Background(), req)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if res.UserID != "1" {
			t.Errorf("Expected response UserID 1, got %s", res.UserID)
		}
	})

	t.Run("Error - Email kosong", func(t *testing.T) {
		repo := &MockUserRepository{
			SaveFunc: func(ctx context.Context, user *User) error {
				t.Fatal("Repository Save() should not be called")
				return nil
			},
		}
		service := &userServiceImpl{repo: repo}

		req := CreateUserRequest{
			UserID: "1",
			Name:   "Erv",
			Email:  "", // Kosong
			Phone:  "08123",
		}

		res, err := service.CreateUser(context.Background(), req)

		if err == nil || err.Error() != "email required" {
			t.Fatalf("Expected 'email required' error, got %v", err)
		}
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
	})

	t.Run("Error - Data tidak lengkap", func(t *testing.T) {
		repo := &MockUserRepository{
			SaveFunc: func(ctx context.Context, user *User) error {
				t.Fatal("Repository Save() should not be called")
				return nil
			},
		}
		service := &userServiceImpl{repo: repo}

		req := CreateUserRequest{
			UserID: "1",
			Name:   "", // Kosong
			Email:  "erv@mail.com",
			Phone:  "08123",
		}

		res, err := service.CreateUser(context.Background(), req)

		if err == nil {
			t.Fatal("Expected validation error, got none")
		}
		if res != nil {
			t.Errorf("Expected nil response, got %v", res)
		}
	})
}

func TestUserService_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success - User ditemukan", func(t *testing.T) {
		repo := &MockUserRepository{
			FindByIDFunc: func(ctx context.Context, userID string) (*User, error) {
				if userID != "1" {
					t.Errorf("Expected UserID 1, got %s", userID)
				}
				return &User{UserID: "1", Name: "Erv"}, nil
			},
		}
		service := &userServiceImpl{repo: repo}

		user, err := service.GetUser(context.Background(), "1")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if user == nil || user.UserID != "1" {
			t.Errorf("Expected User 1, got %v", user)
		}
	})

	t.Run("Error - User tidak ditemukan", func(t *testing.T) {
		repo := &MockUserRepository{
			FindByIDFunc: func(ctx context.Context, userID string) (*User, error) {
				return nil, errors.New("user not found")
			},
		}
		service := &userServiceImpl{repo: repo}

		user, err := service.GetUser(context .Background(), "99")

		if err == nil || err.Error() != "user not found" {
			t.Fatalf("Expected 'user not found' error, got %v", err)
		}
		if user != nil {
			t.Errorf("Expected nil user, got %v", user)
		}
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success - Data berhasil diupdate", func(t *testing.T) {
		repo := &MockUserRepository{
			FindByIDFunc: func(ctx context.Context, userID string) (*User, error) {
				return &User{UserID: "1", Name: "Erv", Email: "erv@mail.com"}, nil
			},
			UpdateFunc: func(ctx context.Context, user *User) error {
				if user.Name != "Erv Update" || user.Email != "erv_update@mail.com" {
					t.Errorf("Updated user data mismatch, got: %+v", user)
				}
				return nil
			},
		}
		service := &userServiceImpl{repo: repo}

		req := UpdateUserRequest{
			Name:  "Erv Update",
			Email: "erv_update@mail.com",
		}

		err := service.UpdateUser(context.Background(), "1", req)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("Error - User tidak ditemukan", func(t *testing.T) {
		repo := &MockUserRepository{
			FindByIDFunc: func(ctx context.Context, userID string) (*User, error) {
				return nil, errors.New("user not found")
			},
		}
		service := &userServiceImpl{repo: repo}

		req := UpdateUserRequest{
			Name:  "Erv Update",
			Email: "erv_update@mail.com",
		}

		err := service.UpdateUser(context.Background(), "99", req)

		if err == nil || err.Error() != "user not found" {
			t.Fatalf("Expected 'user not found' error, got %v", err)
		}
	})
}

func TestUserService_DeactivateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success - User dinonaktifkan", func(t *testing.T) {
		repo := &MockUserRepository{
			FindByIDFunc: func(ctx context.Context, userID string) (*User, error) {
				return &User{UserID: "1", Status: "active"}, nil
			},
			DeactivateFunc: func(ctx context.Context, userID string) error {
				if userID != "1" {
					t.Errorf("Expected UserID 1, got %s", userID)
				}
				return nil
			},
		}
		service := &userServiceImpl{repo: repo}

		err := service.DeactivateUser(context.Background(), "1")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})
}

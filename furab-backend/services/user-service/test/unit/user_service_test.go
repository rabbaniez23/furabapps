package unit

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"furab-backend/services/user-service/internal/model"
	mock_repository "furab-backend/services/user-service/internal/repository/mock"
	"furab-backend/services/user-service/internal/service"
)

// =======================
// 1. CREATE USER
// =======================

func TestCreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)

	mockRepo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(nil)

	svc := service.NewUserService(mockRepo)

	res, err := svc.CreateUser(context.Background(), service.CreateUserRequest{
		UserID: "1",
		Name:   "Erv",
		Email:  "erv@mail.com",
		Phone:  "08123",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.UserID != "1" {
		t.Errorf("expected user_id 1, got %s", res.UserID)
	}
	if res.Message != "sukses" {
		t.Errorf("expected message sukses, got %s", res.Message)
	}
}

func TestCreateUser_EmailEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)

	svc := service.NewUserService(mockRepo)

	_, err := svc.CreateUser(context.Background(), service.CreateUserRequest{
		Name:  "Erv",
		Email: "",
	})

	if err == nil || !errors.Is(err, service.ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestCreateUser_DataTidakLengkap(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)

	svc := service.NewUserService(mockRepo)

	_, err := svc.CreateUser(context.Background(), service.CreateUserRequest{
		Name:  "",
		Email: "erv@mail.com",
	})

	if err == nil || !errors.Is(err, service.ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestCreateUser_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)

	mockRepo.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(errors.New("db error"))

	svc := service.NewUserService(mockRepo)

	_, err := svc.CreateUser(context.Background(), service.CreateUserRequest{
		UserID: "1",
		Name:   "Erv",
		Email:  "erv@mail.com",
	})

	if err == nil || err.Error() != "db error" {
		t.Fatalf("expected db error, got %v", err)
	}
}

// =======================
// 2. GET USER
// =======================

func TestGetUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)

	mockRepo.EXPECT().
		FindByID(gomock.Any(), "1").
		Return(&model.User{UserID: "1", Name: "Erv"}, nil)

	svc := service.NewUserService(mockRepo)

	user, err := svc.GetUser(context.Background(), "1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.UserID != "1" {
		t.Errorf("expected user 1, got %s", user.UserID)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)

	mockRepo.EXPECT().
		FindByID(gomock.Any(), "99").
		Return(nil, nil)

	svc := service.NewUserService(mockRepo)

	_, err := svc.GetUser(context.Background(), "99")

	if err == nil || !errors.Is(err, service.ErrUserNotFound) {
		t.Fatalf("expected user not found, got %v", err)
	}
}

func TestGetUser_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)

	mockRepo.EXPECT().
		FindByID(gomock.Any(), "1").
		Return(nil, errors.New("db error"))

	svc := service.NewUserService(mockRepo)

	_, err := svc.GetUser(context.Background(), "1")

	if err == nil || err.Error() != "db error" {
		t.Fatalf("expected db error, got %v", err)
	}
}

// =======================
// 3. UPDATE USER
// =======================

func TestUpdateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)

	mockRepo.EXPECT().
		FindByID(gomock.Any(), "1").
		Return(&model.User{UserID: "1"}, nil)

	mockRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(nil)

	svc := service.NewUserService(mockRepo)

	err := svc.UpdateUser(context.Background(), "1", service.UpdateUserRequest{
		Name:  "Updated",
		Email: "updated@mail.com",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)

	mockRepo.EXPECT().
		FindByID(gomock.Any(), "99").
		Return(nil, nil)

	svc := service.NewUserService(mockRepo)

	err := svc.UpdateUser(context.Background(), "99", service.UpdateUserRequest{
		Name:  "Updated",
		Email: "updated@mail.com",
	})

	if err == nil || !errors.Is(err, service.ErrUserNotFound) {
		t.Fatalf("expected user not found, got %v", err)
	}
}

func TestUpdateUser_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)

	mockRepo.EXPECT().
		FindByID(gomock.Any(), "1").
		Return(&model.User{UserID: "1"}, nil)

	mockRepo.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Return(errors.New("db error"))

	svc := service.NewUserService(mockRepo)

	err := svc.UpdateUser(context.Background(), "1", service.UpdateUserRequest{
		Name:  "Updated",
		Email: "updated@mail.com",
	})

	if err == nil || err.Error() != "db error" {
		t.Fatalf("expected db error, got %v", err)
	}
}

// =======================
// 4. DEACTIVATE USER
// =======================

func TestDeactivateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)

	mockRepo.EXPECT().
		FindByID(gomock.Any(), "1").
		Return(&model.User{UserID: "1"}, nil)

	mockRepo.EXPECT().
		Deactivate(gomock.Any(), "1").
		Return(nil)

	svc := service.NewUserService(mockRepo)

	err := svc.DeactivateUser(context.Background(), "1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDeactivateUser_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)

	mockRepo.EXPECT().
		FindByID(gomock.Any(), "99").
		Return(nil, nil)

	svc := service.NewUserService(mockRepo)

	err := svc.DeactivateUser(context.Background(), "99")

	if err == nil || !errors.Is(err, service.ErrUserNotFound) {
		t.Fatalf("expected user not found, got %v", err)
	}
}

func TestDeactivateUser_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)

	mockRepo.EXPECT().
		FindByID(gomock.Any(), "1").
		Return(&model.User{UserID: "1"}, nil)

	mockRepo.EXPECT().
		Deactivate(gomock.Any(), "1").
		Return(errors.New("db error"))

	svc := service.NewUserService(mockRepo)

	err := svc.DeactivateUser(context.Background(), "1")

	if err == nil || err.Error() != "db error" {
		t.Fatalf("expected db error, got %v", err)
	}
}
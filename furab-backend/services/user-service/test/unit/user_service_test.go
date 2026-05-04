// Package unit contains unit tests for the user service.
// Unit tests do NOT access any database or external service.
// All dependencies are mocked using gomock.
//
// Tests cover all CRUD operations:
// CreateUser, GetUser, UpdateUser, DeactivateUser
// Each operation is tested for success, error, and edge cases.
package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"furab-backend/services/user-service/internal/model"
	"furab-backend/services/user-service/internal/repository"
	mock_repository "furab-backend/services/user-service/internal/repository/mock"
	"furab-backend/services/user-service/internal/service"

	"go.uber.org/mock/gomock"
)

// --- Helper Functions ---

// newTestService creates a new UserService with mocked dependencies.
func newTestService(t *testing.T) (service.UserService, *mock_repository.MockUserRepository, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	svc := service.NewUserService(mockRepo)
	return svc, mockRepo, ctrl
}

// sampleUser returns a sample User for testing.
func sampleUser() *model.User {
	return &model.User{
		UserID:    "user-123",
		Name:      "John Doe",
		Phone:     "081234567890",
		Email:     "john@example.com",
		Status:    model.UserStatusActive,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

// ========================================
// Test Cases: CreateUser
// ========================================

// TestCreateUser_Success tests creating a user with valid data.
// Expected: user created with active status, response contains user_id and message.
func TestCreateUser_Success(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &model.CreateUserRequest{
		UserID: "user-123",
		Name:   "John Doe",
		Phone:  "081234567890",
		Email:  "john@example.com",
	}

	// Expect repository Create to be called
	mockRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(nil)

	resp, err := svc.CreateUser(ctx, req)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp == nil {
		t.Fatal("expected response, got nil")
	}
	if resp.UserID != "user-123" {
		t.Errorf("expected user_id user-123, got: %s", resp.UserID)
	}
	if resp.Message == "" {
		t.Error("expected non-empty message")
	}
}

// TestCreateUser_NilRequest tests creating a user with nil request.
func TestCreateUser_NilRequest(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.CreateUser(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil request")
	}
	if err != service.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

// TestCreateUser_EmptyUserID tests creating a user with empty user_id.
func TestCreateUser_EmptyUserID(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	req := &model.CreateUserRequest{
		UserID: "",
		Name:   "John Doe",
		Phone:  "081234567890",
		Email:  "john@example.com",
	}

	_, err := svc.CreateUser(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for empty user_id")
	}
}

// TestCreateUser_EmptyName tests creating a user with empty name.
func TestCreateUser_EmptyName(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	req := &model.CreateUserRequest{
		UserID: "user-123",
		Name:   "",
		Phone:  "081234567890",
		Email:  "john@example.com",
	}

	_, err := svc.CreateUser(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

// TestCreateUser_EmptyPhone tests creating a user with empty phone.
func TestCreateUser_EmptyPhone(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	req := &model.CreateUserRequest{
		UserID: "user-123",
		Name:   "John Doe",
		Phone:  "",
		Email:  "john@example.com",
	}

	_, err := svc.CreateUser(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for empty phone")
	}
}

// TestCreateUser_EmptyEmail tests creating a user with empty email.
func TestCreateUser_EmptyEmail(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	req := &model.CreateUserRequest{
		UserID: "user-123",
		Name:   "John Doe",
		Phone:  "081234567890",
		Email:  "",
	}

	_, err := svc.CreateUser(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for empty email")
	}
}

// TestCreateUser_WhitespaceNormalization tests that inputs are trimmed.
func TestCreateUser_WhitespaceNormalization(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	req := &model.CreateUserRequest{
		UserID: "  user-123  ",
		Name:   "  John Doe  ",
		Phone:  "  081234567890  ",
		Email:  "  john@example.com  ",
	}

	mockRepo.EXPECT().
		Create(ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, user *model.User) error {
			if user.UserID != "user-123" {
				t.Errorf("expected trimmed user_id, got: %q", user.UserID)
			}
			if user.Name != "John Doe" {
				t.Errorf("expected trimmed name, got: %q", user.Name)
			}
			if user.Phone != "081234567890" {
				t.Errorf("expected trimmed phone, got: %q", user.Phone)
			}
			if user.Email != "john@example.com" {
				t.Errorf("expected trimmed email, got: %q", user.Email)
			}
			if user.Status != model.UserStatusActive {
				t.Errorf("expected active status, got: %s", user.Status)
			}
			return nil
		})

	resp, err := svc.CreateUser(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.UserID != "user-123" {
		t.Errorf("expected trimmed user_id in response, got: %s", resp.UserID)
	}
}

// TestCreateUser_WhitespaceOnlyFields tests that whitespace-only fields are rejected.
func TestCreateUser_WhitespaceOnlyFields(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	req := &model.CreateUserRequest{
		UserID: "   ",
		Name:   "John",
		Phone:  "08123",
		Email:  "john@example.com",
	}

	_, err := svc.CreateUser(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for whitespace-only user_id")
	}
}

// TestCreateUser_RepositoryError tests creating a user when repository fails.
func TestCreateUser_RepositoryError(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repoErr := errors.New("database connection failed")

	mockRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(repoErr)

	_, err := svc.CreateUser(ctx, &model.CreateUserRequest{
		UserID: "user-123",
		Name:   "John Doe",
		Phone:  "081234567890",
		Email:  "john@example.com",
	})

	if err == nil {
		t.Fatal("expected error from repository")
	}
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error, got: %v", err)
	}
}

// ========================================
// Test Cases: GetUser
// ========================================

// TestGetUser_Success tests retrieving an existing user.
func TestGetUser_Success(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	expected := sampleUser()

	mockRepo.EXPECT().
		GetByID(ctx, expected.UserID).
		Return(expected, nil)

	user, err := svc.GetUser(ctx, expected.UserID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if user.UserID != expected.UserID {
		t.Errorf("expected user_id %s, got: %s", expected.UserID, user.UserID)
	}
	if user.Name != expected.Name {
		t.Errorf("expected name %s, got: %s", expected.Name, user.Name)
	}
	if user.Email != expected.Email {
		t.Errorf("expected email %s, got: %s", expected.Email, user.Email)
	}
	if user.Phone != expected.Phone {
		t.Errorf("expected phone %s, got: %s", expected.Phone, user.Phone)
	}
	if user.Status != model.UserStatusActive {
		t.Errorf("expected status active, got: %s", user.Status)
	}
}

// TestGetUser_NotFound tests retrieving a non-existent user.
func TestGetUser_NotFound(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepo.EXPECT().
		GetByID(ctx, "non-existent").
		Return(nil, repository.ErrUserNotFound)

	_, err := svc.GetUser(ctx, "non-existent")
	if err != service.ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got: %v", err)
	}
}

// TestGetUser_EmptyID tests retrieving a user with empty ID.
func TestGetUser_EmptyID(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.GetUser(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty ID")
	}
	if err != service.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

// TestGetUser_WhitespaceOnlyID tests retrieving a user with whitespace-only ID.
func TestGetUser_WhitespaceOnlyID(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.GetUser(context.Background(), "   ")
	if err == nil {
		t.Fatal("expected error for whitespace-only ID")
	}
	if err != service.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

// TestGetUser_RepositoryError tests retrieving a user when repository fails.
func TestGetUser_RepositoryError(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repoErr := errors.New("database timeout")

	mockRepo.EXPECT().
		GetByID(ctx, "user-123").
		Return(nil, repoErr)

	_, err := svc.GetUser(ctx, "user-123")
	if err == nil {
		t.Fatal("expected error from repository")
	}
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error, got: %v", err)
	}
}

// TestGetUser_TrimmedInput tests that user ID is trimmed before lookup.
func TestGetUser_TrimmedInput(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	expected := sampleUser()

	// Mock expects the trimmed ID
	mockRepo.EXPECT().
		GetByID(ctx, "user-123").
		Return(expected, nil)

	// Pass ID with whitespace
	user, err := svc.GetUser(ctx, "  user-123  ")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if user.UserID != "user-123" {
		t.Errorf("expected user-123, got: %s", user.UserID)
	}
}

// ========================================
// Test Cases: UpdateUser
// ========================================

// TestUpdateUser_Success tests updating an existing user.
func TestUpdateUser_Success(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	existingUser := sampleUser()

	mockRepo.EXPECT().
		GetByID(ctx, existingUser.UserID).
		Return(existingUser, nil)

	mockRepo.EXPECT().
		Update(ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, user *model.User) error {
			if user.Name != "Jane Doe" {
				t.Errorf("expected updated name Jane Doe, got: %s", user.Name)
			}
			if user.Email != "jane@example.com" {
				t.Errorf("expected updated email jane@example.com, got: %s", user.Email)
			}
			if user.Phone != "089876543210" {
				t.Errorf("expected updated phone, got: %s", user.Phone)
			}
			return nil
		})

	resp, err := svc.UpdateUser(ctx, existingUser.UserID, &model.UpdateUserRequest{
		Name:  "Jane Doe",
		Email: "jane@example.com",
		Phone: "089876543210",
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp == nil {
		t.Fatal("expected response, got nil")
	}
	if resp.Message == "" {
		t.Error("expected non-empty message")
	}
}

// TestUpdateUser_NotFound tests updating a non-existent user.
func TestUpdateUser_NotFound(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepo.EXPECT().
		GetByID(ctx, "non-existent").
		Return(nil, repository.ErrUserNotFound)

	_, err := svc.UpdateUser(ctx, "non-existent", &model.UpdateUserRequest{
		Name:  "Jane Doe",
		Email: "jane@example.com",
		Phone: "089876543210",
	})

	if err != service.ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got: %v", err)
	}
}

// TestUpdateUser_EmptyUserID tests updating with empty user ID.
func TestUpdateUser_EmptyUserID(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.UpdateUser(context.Background(), "", &model.UpdateUserRequest{
		Name:  "Jane Doe",
		Email: "jane@example.com",
		Phone: "089876543210",
	})

	if err != service.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

// TestUpdateUser_NilRequest tests updating with nil request.
func TestUpdateUser_NilRequest(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.UpdateUser(context.Background(), "user-123", nil)
	if err != service.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

// TestUpdateUser_EmptyName tests updating with empty name.
func TestUpdateUser_EmptyName(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.UpdateUser(context.Background(), "user-123", &model.UpdateUserRequest{
		Name:  "",
		Email: "jane@example.com",
		Phone: "089876543210",
	})

	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

// TestUpdateUser_EmptyEmail tests updating with empty email.
func TestUpdateUser_EmptyEmail(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.UpdateUser(context.Background(), "user-123", &model.UpdateUserRequest{
		Name:  "Jane Doe",
		Email: "",
		Phone: "089876543210",
	})

	if err == nil {
		t.Fatal("expected error for empty email")
	}
}

// TestUpdateUser_EmptyPhone tests updating with empty phone.
func TestUpdateUser_EmptyPhone(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.UpdateUser(context.Background(), "user-123", &model.UpdateUserRequest{
		Name:  "Jane Doe",
		Email: "jane@example.com",
		Phone: "",
	})

	if err == nil {
		t.Fatal("expected error for empty phone")
	}
}

// TestUpdateUser_RepositoryFindError tests updating when find fails.
func TestUpdateUser_RepositoryFindError(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repoErr := errors.New("database error")

	mockRepo.EXPECT().
		GetByID(ctx, "user-123").
		Return(nil, repoErr)

	_, err := svc.UpdateUser(ctx, "user-123", &model.UpdateUserRequest{
		Name:  "Jane Doe",
		Email: "jane@example.com",
		Phone: "089876543210",
	})

	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error, got: %v", err)
	}
}

// TestUpdateUser_RepositoryUpdateError tests updating when update fails.
func TestUpdateUser_RepositoryUpdateError(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	existingUser := sampleUser()
	repoErr := errors.New("update failed")

	mockRepo.EXPECT().
		GetByID(ctx, existingUser.UserID).
		Return(existingUser, nil)

	mockRepo.EXPECT().
		Update(ctx, gomock.Any()).
		Return(repoErr)

	_, err := svc.UpdateUser(ctx, existingUser.UserID, &model.UpdateUserRequest{
		Name:  "Jane Doe",
		Email: "jane@example.com",
		Phone: "089876543210",
	})

	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error, got: %v", err)
	}
}

// ========================================
// Test Cases: DeactivateUser
// ========================================

// TestDeactivateUser_Success tests deactivating an active user.
// Expected: user status set to inactive, update persisted.
func TestDeactivateUser_Success(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	existingUser := sampleUser()
	existingUser.Status = model.UserStatusActive

	mockRepo.EXPECT().
		GetByID(ctx, existingUser.UserID).
		Return(existingUser, nil)

	mockRepo.EXPECT().
		Update(ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, user *model.User) error {
			if user.Status != model.UserStatusInactive {
				t.Errorf("expected status inactive, got: %s", user.Status)
			}
			return nil
		})

	resp, err := svc.DeactivateUser(ctx, existingUser.UserID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp == nil {
		t.Fatal("expected response, got nil")
	}
	if resp.Message == "" {
		t.Error("expected non-empty message")
	}
}

// TestDeactivateUser_NotFound tests deactivating a non-existent user.
func TestDeactivateUser_NotFound(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRepo.EXPECT().
		GetByID(ctx, "non-existent").
		Return(nil, repository.ErrUserNotFound)

	_, err := svc.DeactivateUser(ctx, "non-existent")
	if err != service.ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got: %v", err)
	}
}

// TestDeactivateUser_EmptyUserID tests deactivating with empty user ID.
func TestDeactivateUser_EmptyUserID(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.DeactivateUser(context.Background(), "")
	if err != service.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

// TestDeactivateUser_WhitespaceOnlyID tests deactivating with whitespace-only ID.
func TestDeactivateUser_WhitespaceOnlyID(t *testing.T) {
	svc, _, ctrl := newTestService(t)
	defer ctrl.Finish()

	_, err := svc.DeactivateUser(context.Background(), "   ")
	if err != service.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
}

// TestDeactivateUser_RepositoryFindError tests deactivating when find fails.
func TestDeactivateUser_RepositoryFindError(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repoErr := errors.New("database connection lost")

	mockRepo.EXPECT().
		GetByID(ctx, "user-123").
		Return(nil, repoErr)

	_, err := svc.DeactivateUser(ctx, "user-123")
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error, got: %v", err)
	}
}

// TestDeactivateUser_RepositoryUpdateError tests deactivating when update fails.
func TestDeactivateUser_RepositoryUpdateError(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	existingUser := sampleUser()
	repoErr := errors.New("update failed")

	mockRepo.EXPECT().
		GetByID(ctx, existingUser.UserID).
		Return(existingUser, nil)

	mockRepo.EXPECT().
		Update(ctx, gomock.Any()).
		Return(repoErr)

	_, err := svc.DeactivateUser(ctx, existingUser.UserID)
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error, got: %v", err)
	}
}

// TestDeactivateUser_TrimmedInput tests that user ID is trimmed before deactivation.
func TestDeactivateUser_TrimmedInput(t *testing.T) {
	svc, mockRepo, ctrl := newTestService(t)
	defer ctrl.Finish()

	ctx := context.Background()
	existingUser := sampleUser()

	// Mock expects the trimmed ID
	mockRepo.EXPECT().
		GetByID(ctx, "user-123").
		Return(existingUser, nil)

	mockRepo.EXPECT().
		Update(ctx, gomock.Any()).
		Return(nil)

	// Pass ID with whitespace
	resp, err := svc.DeactivateUser(ctx, "  user-123  ")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp == nil {
		t.Fatal("expected response, got nil")
	}
}

// ========================================
// Test Cases: UserStatus Validation
// ========================================

// TestUserStatus_IsValid tests the status validation.
func TestUserStatus_IsValid(t *testing.T) {
	tests := []struct {
		status model.UserStatus
		valid  bool
	}{
		{model.UserStatusActive, true},
		{model.UserStatusInactive, true},
		{model.UserStatus("unknown"), false},
		{model.UserStatus(""), false},
	}

	for _, tc := range tests {
		result := tc.status.IsValid()
		if result != tc.valid {
			t.Errorf("status %q: expected IsValid=%v, got %v", tc.status, tc.valid, result)
		}
	}
}

// ========================================
// Test Cases: Request Validation (Model)
// ========================================

// TestCreateUserRequest_Validate tests the CreateUserRequest validation.
func TestCreateUserRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     model.CreateUserRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: model.CreateUserRequest{
				UserID: "user-1",
				Name:   "John",
				Phone:  "08123",
				Email:  "john@example.com",
			},
			wantErr: false,
		},
		{
			name: "empty user_id",
			req: model.CreateUserRequest{
				Name:  "John",
				Phone: "08123",
				Email: "john@example.com",
			},
			wantErr: true,
		},
		{
			name: "empty name",
			req: model.CreateUserRequest{
				UserID: "user-1",
				Phone:  "08123",
				Email:  "john@example.com",
			},
			wantErr: true,
		},
		{
			name: "empty phone",
			req: model.CreateUserRequest{
				UserID: "user-1",
				Name:   "John",
				Email:  "john@example.com",
			},
			wantErr: true,
		},
		{
			name: "empty email",
			req: model.CreateUserRequest{
				UserID: "user-1",
				Name:   "John",
				Phone:  "08123",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.req.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

// TestUpdateUserRequest_Validate tests the UpdateUserRequest validation.
func TestUpdateUserRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     model.UpdateUserRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: model.UpdateUserRequest{
				Name:  "Jane",
				Phone: "08123",
				Email: "jane@example.com",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			req: model.UpdateUserRequest{
				Phone: "08123",
				Email: "jane@example.com",
			},
			wantErr: true,
		},
		{
			name: "empty phone",
			req: model.UpdateUserRequest{
				Name:  "Jane",
				Email: "jane@example.com",
			},
			wantErr: true,
		},
		{
			name: "empty email",
			req: model.UpdateUserRequest{
				Name:  "Jane",
				Phone: "08123",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.req.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

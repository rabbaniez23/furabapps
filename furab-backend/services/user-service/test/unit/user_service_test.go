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

var (
	errRepoSave       = errors.New("repo save error")
	errRepoFindByID   = errors.New("repo find error")
	errRepoUpdate     = errors.New("repo update error")
	errRepoDeactivate = errors.New("repo deactivate error")
)

type userArgMatcher struct {
	match func(*model.User) bool
}

func (m userArgMatcher) Matches(x any) bool {
	u, ok := x.(*model.User)
	if !ok {
		return false
	}
	return m.match(u)
}

func (m userArgMatcher) String() string {
	return "matches *model.User predicate"
}

func matchUser(match func(*model.User) bool) gomock.Matcher {
	return userArgMatcher{match: match}
}

func setupUserService(t *testing.T) (*gomock.Controller, *mock_repository.MockUserRepository, service.UserService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	return ctrl, mockRepo, service.NewUserService(mockRepo)
}

func TestUserService_CreateUser(t *testing.T) {
	ctrl, mockRepo, svc := setupUserService(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().
			Save(gomock.Any(), matchUser(func(u *model.User) bool {
				return u != nil &&
					u.UserID == "user-1" &&
					u.Name == "Erv" &&
					u.Email == "erv@mail.com" &&
					u.Phone == "08123"
			})).
			Return(nil)

		res, err := svc.CreateUser(context.Background(), service.CreateUserRequest{
			UserID: " user-1 ",
			Name:   " Erv ",
			Email:  " erv@mail.com ",
			Phone:  " 08123 ",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if res == nil || res.UserID != "user-1" {
			t.Fatalf("expected user-1, got %#v", res)
		}
		if res.Message != "user created" {
			t.Fatalf("expected message 'user created', got %q", res.Message)
		}
	})

	t.Run("validation_error_user_id_required", func(t *testing.T) {
		res, err := svc.CreateUser(context.Background(), service.CreateUserRequest{
			UserID: "",
			Name:   "Erv",
			Email:  "erv@mail.com",
			Phone:  "08123",
		})
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, service.ErrUserIDRequired) {
			t.Fatalf("expected ErrUserIDRequired, got %v", err)
		}
	})

	t.Run("validation_error_required_fields", func(t *testing.T) {
		_, errName := svc.CreateUser(context.Background(), service.CreateUserRequest{
			UserID: "user-1",
			Name:   "",
			Email:  "erv@mail.com",
			Phone:  "08123",
		})
		if !errors.Is(errName, service.ErrNameRequired) {
			t.Fatalf("expected ErrNameRequired, got %v", errName)
		}

		_, errEmail := svc.CreateUser(context.Background(), service.CreateUserRequest{
			UserID: "user-1",
			Name:   "Erv",
			Email:  "",
			Phone:  "08123",
		})
		if !errors.Is(errEmail, service.ErrEmailRequired) {
			t.Fatalf("expected ErrEmailRequired, got %v", errEmail)
		}

		_, errPhone := svc.CreateUser(context.Background(), service.CreateUserRequest{
			UserID: "user-1",
			Name:   "Erv",
			Email:  "erv@mail.com",
			Phone:  "",
		})
		if !errors.Is(errPhone, service.ErrPhoneRequired) {
			t.Fatalf("expected ErrPhoneRequired, got %v", errPhone)
		}
	})

	t.Run("repository_error", func(t *testing.T) {
		mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(errRepoSave)

		res, err := svc.CreateUser(context.Background(), service.CreateUserRequest{
			UserID: "user-1",
			Name:   "Erv",
			Email:  "erv@mail.com",
			Phone:  "08123",
		})
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, errRepoSave) {
			t.Fatalf("expected repo save error, got %v", err)
		}
	})

	t.Run("context_cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		res, err := svc.CreateUser(ctx, service.CreateUserRequest{
			UserID: "user-1",
			Name:   "Erv",
			Email:  "erv@mail.com",
			Phone:  "08123",
		})
		if res != nil {
			t.Fatalf("expected nil response, got %#v", res)
		}
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	})
}

func TestUserService_GetUser(t *testing.T) {
	ctrl, mockRepo, svc := setupUserService(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(gomock.Any(), "1").Return(&model.User{UserID: "1", Name: "Erv"}, nil)

		user, err := svc.GetUser(context.Background(), "1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if user == nil || user.UserID != "1" {
			t.Fatalf("expected user 1, got %#v", user)
		}
	})

	t.Run("validation_error_user_id_required", func(t *testing.T) {
		user, err := svc.GetUser(context.Background(), " ")
		if user != nil {
			t.Fatalf("expected nil user, got %#v", user)
		}
		if !errors.Is(err, service.ErrUserIDRequired) {
			t.Fatalf("expected ErrUserIDRequired, got %v", err)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(gomock.Any(), "99").Return(nil, nil)

		user, err := svc.GetUser(context.Background(), "99")
		if user != nil {
			t.Fatalf("expected nil user, got %#v", user)
		}
		if !errors.Is(err, service.ErrUserNotFound) {
			t.Fatalf("expected ErrUserNotFound, got %v", err)
		}
	})

	t.Run("repository_error_find", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(gomock.Any(), "1").Return(nil, errRepoFindByID)

		user, err := svc.GetUser(context.Background(), "1")
		if user != nil {
			t.Fatalf("expected nil user, got %#v", user)
		}
		if !errors.Is(err, errRepoFindByID) {
			t.Fatalf("expected repo find error, got %v", err)
		}
	})

	t.Run("context_cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		user, err := svc.GetUser(ctx, "1")
		if user != nil {
			t.Fatalf("expected nil user, got %#v", user)
		}
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	ctrl, mockRepo, svc := setupUserService(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(gomock.Any(), "1").Return(&model.User{UserID: "1"}, nil)
		mockRepo.EXPECT().
			Update(gomock.Any(), matchUser(func(u *model.User) bool {
				return u != nil && u.UserID == "1" && u.Name == "Updated" && u.Email == "updated@mail.com"
			})).
			Return(nil)

		err := svc.UpdateUser(context.Background(), "1", service.UpdateUserRequest{
			Name:  "Updated",
			Email: "updated@mail.com",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("validation_error_user_id_required", func(t *testing.T) {
		err := svc.UpdateUser(context.Background(), "", service.UpdateUserRequest{
			Name:  "Updated",
			Email: "updated@mail.com",
		})
		if !errors.Is(err, service.ErrUserIDRequired) {
			t.Fatalf("expected ErrUserIDRequired, got %v", err)
		}
	})

	t.Run("validation_error_payload", func(t *testing.T) {
		err := svc.UpdateUser(context.Background(), "1", service.UpdateUserRequest{
			Name:  "",
			Email: "updated@mail.com",
		})
		if !errors.Is(err, service.ErrNameRequired) {
			t.Fatalf("expected ErrNameRequired, got %v", err)
		}
	})

	t.Run("validation_error_payload_email_required", func(t *testing.T) {
		err := svc.UpdateUser(context.Background(), "1", service.UpdateUserRequest{
			Name:  "Updated",
			Email: " ",
		})
		if !errors.Is(err, service.ErrEmailRequired) {
			t.Fatalf("expected ErrEmailRequired, got %v", err)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(gomock.Any(), "99").Return(nil, nil)

		err := svc.UpdateUser(context.Background(), "99", service.UpdateUserRequest{
			Name:  "Updated",
			Email: "updated@mail.com",
		})
		if !errors.Is(err, service.ErrUserNotFound) {
			t.Fatalf("expected ErrUserNotFound, got %v", err)
		}
	})

	t.Run("repository_error_find", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(gomock.Any(), "1").Return(nil, errRepoFindByID)

		err := svc.UpdateUser(context.Background(), "1", service.UpdateUserRequest{
			Name:  "Updated",
			Email: "updated@mail.com",
		})
		if !errors.Is(err, errRepoFindByID) {
			t.Fatalf("expected repo find error, got %v", err)
		}
	})

	t.Run("repository_error_update", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(gomock.Any(), "1").Return(&model.User{UserID: "1"}, nil)
		mockRepo.EXPECT().
			Update(gomock.Any(), matchUser(func(u *model.User) bool {
				return u != nil && u.UserID == "1" && u.Name == "Updated" && u.Email == "updated@mail.com"
			})).
			Return(errRepoUpdate)

		err := svc.UpdateUser(context.Background(), "1", service.UpdateUserRequest{
			Name:  "Updated",
			Email: "updated@mail.com",
		})
		if !errors.Is(err, errRepoUpdate) {
			t.Fatalf("expected repo update error, got %v", err)
		}
	})

	t.Run("context_cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := svc.UpdateUser(ctx, "1", service.UpdateUserRequest{
			Name:  "Updated",
			Email: "updated@mail.com",
		})
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	})
}

func TestUserService_DeactivateUser(t *testing.T) {
	ctrl, mockRepo, svc := setupUserService(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(gomock.Any(), "1").Return(&model.User{UserID: "1"}, nil)
		mockRepo.EXPECT().Deactivate(gomock.Any(), "1").Return(nil)

		err := svc.DeactivateUser(context.Background(), "1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("validation_error_user_id_required", func(t *testing.T) {
		err := svc.DeactivateUser(context.Background(), "")
		if !errors.Is(err, service.ErrUserIDRequired) {
			t.Fatalf("expected ErrUserIDRequired, got %v", err)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(gomock.Any(), "99").Return(nil, nil)

		err := svc.DeactivateUser(context.Background(), "99")
		if !errors.Is(err, service.ErrUserNotFound) {
			t.Fatalf("expected ErrUserNotFound, got %v", err)
		}
	})

	t.Run("repository_error_find", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(gomock.Any(), "1").Return(nil, errRepoFindByID)

		err := svc.DeactivateUser(context.Background(), "1")
		if !errors.Is(err, errRepoFindByID) {
			t.Fatalf("expected repo find error, got %v", err)
		}
	})

	t.Run("repository_error_deactivate", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(gomock.Any(), "1").Return(&model.User{UserID: "1"}, nil)
		mockRepo.EXPECT().Deactivate(gomock.Any(), "1").Return(errRepoDeactivate)

		err := svc.DeactivateUser(context.Background(), "1")
		if !errors.Is(err, errRepoDeactivate) {
			t.Fatalf("expected repo deactivate error, got %v", err)
		}
	})

	t.Run("normalize_user_id_before_repo_call", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(gomock.Any(), "1").Return(&model.User{UserID: "1"}, nil)
		mockRepo.EXPECT().Deactivate(gomock.Any(), "1").Return(nil)

		err := svc.DeactivateUser(context.Background(), " 1 ")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("context_cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := svc.DeactivateUser(ctx, "1")
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	})
}

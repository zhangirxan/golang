package service

import (
	"errors"
	"practice-8/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// ─── GetUserByID ──────────────────────────────────────────────────────────────

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Bakytzhan Agai"}
	mockRepo.EXPECT().GetUserByID(1).Return(user, nil)

	result, err := svc.GetUserByID(1)
	require.NoError(t, err)
	assert.Equal(t, user, result)
}

// ─── CreateUser ───────────────────────────────────────────────────────────────

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	user := &repository.User{ID: 2, Name: "Test User"}
	mockRepo.EXPECT().CreateUser(user).Return(nil)

	err := svc.CreateUser(user)
	require.NoError(t, err)
}

// ─── RegisterUser ─────────────────────────────────────────────────────────────

func TestRegisterUser_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	existing := &repository.User{ID: 5, Name: "Existing", Email: "test@mail.com"}
	mockRepo.EXPECT().GetByEmail("test@mail.com").Return(existing, nil)

	newUser := &repository.User{Name: "New Guy"}
	err := svc.RegisterUser(newUser, "test@mail.com")
	require.Error(t, err)
	assert.Equal(t, "user with this email already exists", err.Error())
}

func TestRegisterUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	newUser := &repository.User{ID: 3, Name: "New User", Email: "new@mail.com"}
	mockRepo.EXPECT().GetByEmail("new@mail.com").Return(nil, nil)
	mockRepo.EXPECT().CreateUser(newUser).Return(nil)

	err := svc.RegisterUser(newUser, "new@mail.com")
	require.NoError(t, err)
}

func TestRegisterUser_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	mockRepo.EXPECT().GetByEmail("err@mail.com").Return(nil, errors.New("db error"))

	newUser := &repository.User{Name: "Someone"}
	err := svc.RegisterUser(newUser, "err@mail.com")
	require.Error(t, err)
	assert.Equal(t, "error getting user with this email", err.Error())
}

// ─── UpdateUserName ───────────────────────────────────────────────────────────

func TestUpdateUserName_EmptyName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	err := svc.UpdateUserName(1, "")
	require.Error(t, err)
	assert.Equal(t, "name cannot be empty", err.Error())
}

func TestUpdateUserName_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	mockRepo.EXPECT().GetUserByID(99).Return(nil, errors.New("user not found"))

	err := svc.UpdateUserName(99, "New Name")
	require.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}

func TestUpdateUserName_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	user := &repository.User{ID: 2, Name: "Old Name"}
	mockRepo.EXPECT().GetUserByID(2).Return(user, nil)
	mockRepo.EXPECT().UpdateUser(user).Return(nil)

	err := svc.UpdateUserName(2, "New Name")
	require.NoError(t, err)

	// Verify name was actually changed before update
	assert.Equal(t, "New Name", user.Name)
}

func TestUpdateUserName_UpdateFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	user := &repository.User{ID: 3, Name: "Old"}
	mockRepo.EXPECT().GetUserByID(3).Return(user, nil)
	mockRepo.EXPECT().UpdateUser(user).Return(errors.New("update failed"))

	err := svc.UpdateUserName(3, "New")
	require.Error(t, err)
	assert.Equal(t, "update failed", err.Error())
}

// ─── DeleteUser ───────────────────────────────────────────────────────────────

func TestDeleteUser_Admin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	err := svc.DeleteUser(1)
	require.Error(t, err)
	assert.Equal(t, "it is not allowed to delete admin user", err.Error())
}

func TestDeleteUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	mockRepo.EXPECT().DeleteUser(5).Return(nil)

	err := svc.DeleteUser(5)
	require.NoError(t, err)
}

func TestDeleteUser_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	mockRepo.EXPECT().DeleteUser(7).Return(errors.New("db error"))

	err := svc.DeleteUser(7)
	require.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}
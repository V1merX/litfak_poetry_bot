package services

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/V1merX/litfak_poetry_bot/internal/domain"
	"github.com/V1merX/litfak_poetry_bot/internal/services/mocks"
	"github.com/gookit/goutil/testutil/assert"
	"go.uber.org/mock/gomock"
)

func TestUserService_GetByTelegramID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name       string
		telegramID int64
		user       *domain.User
		err        error
	}{
		{
			name:       "no_error",
			telegramID: 1234,
			user: &domain.User{
				UserID:     123,
				TelegramID: 1234,
				ChatID:     1234,
				UserName:   "test_username",
				FirstName:  "test_firstname",
				LastName:   "test_lastname",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			err: nil,
		},
		{
			name:       "with_error",
			telegramID: 1234,
			user:       nil,
			err:        errors.New("failed to get user by telegram id"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockRepo := mocks.NewMockUserRepository(ctrl)
			mockRepo.EXPECT().GetByTelegramID(gomock.Any(), c.telegramID).Return(c.user, c.err).Times(1)

			userService := NewUserService(slog.Default(), mockRepo)
			user, err := userService.GetByTelegramID(context.Background(), c.telegramID)

			assert.Equal(t, c.err, err)
			assert.Equal(t, c.user, user)
		})
	}
}

func TestUserService_GetByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name   string
		userID int64
		user   *domain.User
		err    error
	}{
		{
			name:   "no_error",
			userID: 1234,
			user: &domain.User{
				UserID:     1234,
				TelegramID: 12345,
				ChatID:     1234,
				UserName:   "test_username",
				FirstName:  "test_firstname",
				LastName:   "test_lastname",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			err: nil,
		},
		{
			name:   "with_error",
			userID: 1234,
			user:   nil,
			err:    errors.New("failed to get user by id"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockRepo := mocks.NewMockUserRepository(ctrl)
			mockRepo.EXPECT().GetByUserID(gomock.Any(), c.userID).Return(c.user, c.err).Times(1)

			userService := NewUserService(slog.Default(), mockRepo)
			user, err := userService.GetByUserID(context.Background(), c.userID)

			assert.Equal(t, c.err, err)
			assert.Equal(t, c.user, user)
		})
	}
}

func TestUserService_NewUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name           string
		user           *domain.User
		expectedUserID int64
		err            error
	}{
		{
			name: "no_error",
			user: &domain.User{
				TelegramID: 12345,
				ChatID:     123456,
				UserName:   "test_username",
				FirstName:  "test_firstname",
				LastName:   "test_lastname",
			},
			expectedUserID: 1,
			err:            nil,
		},
		{
			name: "with_error",
			user: &domain.User{
				TelegramID: 12345,
				ChatID:     123456,
				UserName:   "test_username",
				FirstName:  "test_firstname",
				LastName:   "test_lastname",
			},
			expectedUserID: 0,
			err:            errors.New("failed to create a new user"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockRepo := mocks.NewMockUserRepository(ctrl)
			mockRepo.EXPECT().NewUser(gomock.Any(), c.user).Return(c.expectedUserID, c.err).Times(1)

			userService := NewUserService(slog.Default(), mockRepo)
			userID, err := userService.NewUser(context.Background(), c.user)

			assert.Equal(t, c.err, err)
			assert.Equal(t, c.expectedUserID, userID)
		})
	}
}

func TestUserService_GetAllChatIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name    string
		chatIDs *[]domain.User
		err     error
	}{
		{
			name: "no_error",
			chatIDs: &[]domain.User{
				{
					UserID:     1,
					TelegramID: 12345,
					ChatID:     123456,
					UserName:   "test1_username",
					FirstName:  "test1_firstname",
					LastName:   "test1_lastname",
				},
				{
					UserID:     2,
					TelegramID: 123456,
					ChatID:     1234567,
					UserName:   "test2_username",
					FirstName:  "test2_firstname",
					LastName:   "test2_lastname",
				},
				{
					UserID:     3,
					TelegramID: 1234567,
					ChatID:     12345678,
					UserName:   "test3_username",
					FirstName:  "test3_firstname",
					LastName:   "test3_lastname",
				},
			},
			err: nil,
		},
		{
			name:    "with_error",
			chatIDs: nil,
			err:     errors.New("failed to get all chat IDs"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockRepo := mocks.NewMockUserRepository(ctrl)
			mockRepo.EXPECT().GetAllChatIDs(gomock.Any()).Return(c.chatIDs, c.err).Times(1)

			userService := NewUserService(slog.Default(), mockRepo)
			chatIDs, err := userService.GetAllChatIDs(context.Background())

			assert.Equal(t, c.err, err)
			assert.Equal(t, c.chatIDs, chatIDs)
		})
	}
}

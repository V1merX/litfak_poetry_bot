package services

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/V1merX/litfak_poetry_bot/internal/domain"
	"github.com/V1merX/litfak_poetry_bot/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPoemService_GetActualPoem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name         string
		expectedPoem *domain.Poem
		err          error
	}{
		{
			name: "no_error",
			expectedPoem: &domain.Poem{
				PoemID:    1,
				Name:      "test",
				AuthorID:  3,
				Text:      "test_text",
				Meta:      "test_meta",
				IsSent:    false,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			err: nil,
		},
		{
			name:         "with_error",
			expectedPoem: nil,
			err:          errors.New("failed to get actual poem"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockRepo := mocks.NewMockPoemRepository(ctrl)
			mockRepo.EXPECT().
				GetActualPoem(gomock.Any()).
				Return(c.expectedPoem, c.err).
				Times(1)

			service := NewPoemService(slog.Default(), mockRepo)

			poem, err := service.GetActualPoem(context.Background())

			assert.Equal(t, c.err, err)
			assert.Equal(t, c.expectedPoem, poem)
		})
	}
}

func TestPoemService_UpdateStatusSentPoem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name   string
		poemID int64
		isSent bool
		err    error
	}{
		{
			name:   "no_error",
			poemID: 2,
			isSent: true,
			err:    nil,
		},
		{
			name:   "with_error",
			poemID: 2,
			isSent: false,
			err:    errors.New("failed to update status sent poem"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockRepo := mocks.NewMockPoemRepository(ctrl)
			mockRepo.EXPECT().
				UpdateStatusSentPoem(gomock.Any(), c.poemID, c.isSent).
				Return(c.err).
				Times(1)

			service := NewPoemService(slog.Default(), mockRepo)

			err := service.UpdateStatusSentPoem(context.Background(), c.poemID, c.isSent)

			assert.Equal(t, c.err, err)
		})
	}
}
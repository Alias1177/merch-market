package info_test

import (
	"context"
	"errors"
	"github.com/Alias1177/merch-store/internal/models"
	"github.com/Alias1177/merch-store/internal/usecase/info"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockInfoRepository struct {
	mock.Mock
}

func (m *MockInfoRepository) GetUserInfo(ctx context.Context, userID int) (*models.InfoResponse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*models.InfoResponse), args.Error(1)
}

func TestInfoUsecase_GetUserInfo(t *testing.T) {
	mockRepo := new(MockInfoRepository)
	usecase := info.NewInfoUsecase(mockRepo)

	tests := []struct {
		name       string
		userID     int
		mockResult *models.InfoResponse
		mockError  error
		wantResult *models.InfoResponse
		wantErr    bool
	}{
		{
			name:       "success",
			userID:     1,
			mockResult: &models.InfoResponse{},
			mockError:  nil,
			wantResult: &models.InfoResponse{},
			wantErr:    false,
		},
		{
			name:       "repository error",
			userID:     2,
			mockResult: nil,
			mockError:  errors.New("repository error"),
			wantResult: nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("GetUserInfo", mock.Anything, tt.userID).Return(tt.mockResult, tt.mockError)

			result, err := usecase.GetUserInfo(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantResult, result)
			mockRepo.AssertExpectations(t)
		})
	}
}

package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Alias1177/merch-store/internal/models"
	"github.com/Alias1177/merch-store/internal/usecase/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDBRepo struct {
	mock.Mock
}

func (m *MockDBRepo) CreateUser(ctx context.Context, username, passwordHash string, coins int) (*models.User, error) {
	args := m.Called(ctx, username, passwordHash, coins)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestUserUsecase_CreateUser(t *testing.T) {
	mockRepo := new(MockDBRepo)
	usecase := auth.New(mockRepo, "secret")

	tests := []struct {
		name       string
		reqData    models.RegisterRequest
		mockUser   *models.User
		mockError  error
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:    "successful user creation",
			reqData: models.RegisterRequest{Username: "user1", Password: "password123"},
			mockUser: &models.User{
				ID:       1,
				Username: "user1",
			},
			mockError: nil,
			wantErr:   false,
		},
		{
			name:       "user already exists",
			reqData:    models.RegisterRequest{Username: "user1", Password: "password123"},
			mockUser:   nil,
			mockError:  errors.New("user already exists"),
			wantErr:    true,
			wantErrMsg: "user already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil

			mockRepo.On("CreateUser", mock.Anything, tt.reqData.Username, mock.Anything, 1000).
				Return(tt.mockUser, tt.mockError)

			_, err := usecase.CreateUser(context.Background(), tt.reqData)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrMsg != "" {
					assert.Contains(t, err.Error(), tt.wantErrMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

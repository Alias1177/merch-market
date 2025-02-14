package coins_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Alias1177/merch-store/internal/usecase/coins"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCoinsRepository struct {
	mock.Mock
}

func (m *MockCoinsRepository) SendCoins(ctx context.Context, senderID int, receiverUsername string, amount int) error {
	args := m.Called(ctx, senderID, receiverUsername, amount)
	return args.Error(0)
}

func TestCoinsUsecase_SendCoins(t *testing.T) {
	tests := []struct {
		name      string
		senderID  int
		receiver  string
		amount    int
		mockError error
		wantErr   bool
	}{
		{
			name:      "successful transaction",
			senderID:  1,
			receiver:  "receiver1",
			amount:    500,
			mockError: nil,
			wantErr:   false,
		},
		{
			name:      "invalid amount",
			senderID:  1,
			receiver:  "receiver1",
			amount:    -10,
			mockError: nil,
			wantErr:   true,
		},
		{
			name:      "repository error",
			senderID:  1,
			receiver:  "receiver1",
			amount:    500,
			mockError: errors.New("repository error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем новый мок для каждого тест-кейса
			mockRepo := new(MockCoinsRepository)
			usecase := coins.NewCoinsUsecase(mockRepo)

			// Настраиваем мок только если ожидается вызов репозитория
			if tt.amount > 0 {
				mockRepo.On("SendCoins", mock.Anything, tt.senderID, tt.receiver, tt.amount).Return(tt.mockError)
			}

			err := usecase.SendCoins(context.Background(), tt.senderID, tt.receiver, tt.amount)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

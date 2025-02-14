package buy_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Alias1177/merch-store/internal/usecase/buy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBuyRepo struct {
	mock.Mock
}

func (m *MockBuyRepo) BuyItem(ctx context.Context, userID, itemID int) error {
	args := m.Called(ctx, userID, itemID)
	return args.Error(0)
}

func TestBuyUsecase_BuyItem(t *testing.T) {
	tests := []struct {
		name      string
		userID    int
		itemID    int
		mockError error
		wantErr   bool
	}{
		{
			name:      "successful purchase",
			userID:    1,
			itemID:    100,
			mockError: nil,
			wantErr:   false,
		},
		{
			name:      "repository error",
			userID:    1,
			itemID:    100,
			mockError: errors.New("repository error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем новый мок для каждого тест-кейса
			mockRepo := new(MockBuyRepo)
			usecase := buy.NewBuyUsecase(mockRepo)

			mockRepo.On("BuyItem", mock.Anything, tt.userID, tt.itemID).Return(tt.mockError)

			err := usecase.BuyItem(context.Background(), tt.userID, tt.itemID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

package handlers

import (
	"context"
	"encoding/json"
	"github.com/Alias1177/merch-store/internal/constants"
	"golang.org/x/crypto/bcrypt"

	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Alias1177/merch-store/internal/models"
	"github.com/Alias1177/merch-store/internal/usecase/auth"
	"github.com/Alias1177/merch-store/internal/usecase/buy"
	"github.com/Alias1177/merch-store/internal/usecase/coins"
	"github.com/Alias1177/merch-store/internal/usecase/info"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDBRepo - мок-реализация репозитория базы данных
type MockDBRepo struct {
	mock.Mock
}

func (m *MockDBRepo) CreateUser(ctx context.Context, username, passwordHash string, coins int) (*models.User, error) {
	args := m.Called(ctx, username, passwordHash, coins)
	user, _ := args.Get(0).(*models.User)
	return user, args.Error(1)
}

func (m *MockDBRepo) BuyItem(ctx context.Context, userID, itemID int) error {
	args := m.Called(ctx, userID, itemID)
	return args.Error(0)
}

func (m *MockDBRepo) SendCoins(ctx context.Context, senderID int, receiverUsername string, amount int) error {
	args := m.Called(ctx, senderID, receiverUsername, amount)
	return args.Error(0)
}

func (m *MockDBRepo) GetUserInfo(ctx context.Context, userID int) (*models.InfoResponse, error) {
	args := m.Called(ctx, userID)
	info, _ := args.Get(0).(*models.InfoResponse)
	return info, args.Error(1)
}

func TestCreateUser(t *testing.T) {
	mockRepo := new(MockDBRepo)
	userUsecase := auth.New(mockRepo, "mockSecret")

	// Используем mock.MatchedBy для проверки пароля
	mockRepo.On("CreateUser", mock.Anything, "testuser", mock.MatchedBy(func(hashedPassword string) bool {
		// Проверяем, является ли строка действительным bcrypt-хэшем
		return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte("password123")) == nil
	}), 1000).Return(&models.User{
		ID:           1,
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Coins:        1000,
	}, nil)

	token, err := userUsecase.CreateUser(context.Background(), models.RegisterRequest{
		Username: "testuser",
		Password: "password123",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

// Тест для покупки предмета
func TestBuyItem(t *testing.T) {
	mockRepo := new(MockDBRepo)
	buyUsecase := buy.NewBuyUsecase(mockRepo)

	mockRepo.On("BuyItem", mock.Anything, 1, 2).Return(nil)

	err := buyUsecase.BuyItem(context.Background(), 1, 2)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Тест на отправку монет
func TestSendCoins(t *testing.T) {
	mockRepo := new(MockDBRepo)
	coinsUsecase := coins.NewCoinsUsecase(mockRepo)

	mockRepo.On("SendCoins", mock.Anything, 1, "receiver", 100).Return(nil)

	err := coinsUsecase.SendCoins(context.Background(), 1, "receiver", 100)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Тест для получения информации
func TestGetUserInfo(t *testing.T) {
	mockRepo := new(MockDBRepo)
	infoUsecase := info.NewInfoUsecase(mockRepo)

	expectedInfo := &models.InfoResponse{
		Coins: 1000,
		Inventory: []models.InventoryItem{
			{Type: "item1", Quantity: 2},
		},
		CoinHistory: models.CoinHistoryDetails{
			Received: []models.ReceivedTransaction{
				{FromUser: "user1", Amount: 100},
			},
			Sent: []models.SentTransaction{
				{ToUser: "user2", Amount: 50},
			},
		},
	}

	mockRepo.On("GetUserInfo", mock.Anything, 1).Return(expectedInfo, nil)

	info, err := infoUsecase.GetUserInfo(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedInfo, info)
	mockRepo.AssertExpectations(t)
}

// Тест для обработчика регистрации
func TestRegisterHandler(t *testing.T) {
	mockRepo := new(MockDBRepo)
	userUsecase := auth.New(mockRepo, "mockSecret")
	handler := New(userUsecase, nil, nil, nil)

	mockRepo.On("CreateUser", mock.Anything, "testuser", mock.Anything, 1000).Return(&models.User{
		ID:           1,
		Username:     "testuser",
		PasswordHash: "hashedpassword",
		Coins:        1000,
	}, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/auth", strings.NewReader(`{"username":"testuser","password":"password123"}`))
	req.Header.Add("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.RegisterHandler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestHandleInfo(t *testing.T) {
	mockRepo := new(MockDBRepo)
	infoUsecase := info.NewInfoUsecase(mockRepo)
	handler := New(nil, nil, infoUsecase, nil)

	expectedInfo := &models.InfoResponse{
		Coins: 1000,
		Inventory: []models.InventoryItem{
			{Type: "item1", Quantity: 2},
		},
		CoinHistory: models.CoinHistoryDetails{
			Received: []models.ReceivedTransaction{
				{FromUser: "user1", Amount: 100},
			},
			Sent: []models.SentTransaction{
				{ToUser: "user2", Amount: 50},
			},
		},
	}

	// Настроим mock для метода `GetUserInfo`
	mockRepo.On("GetUserInfo", mock.MatchedBy(func(ctx context.Context) bool {
		// Достаем значение из контекста и проверяем его
		val, ok := ctx.Value(constants.UserIDContextKey).(int)
		return ok && val == 1 // Проверяем, что значение - это 1
	}), 1).Return(expectedInfo, nil)

	// Создаем реквест с контекстом, который содержит userID = 1
	req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
	req = req.WithContext(context.WithValue(req.Context(), constants.UserIDContextKey, 1)) // Передаем userID = 1

	// Создаем recorder для получения ответа
	rec := httptest.NewRecorder()

	// Запускаем хендлер
	handler.HandleInfo(rec, req)

	// Проверяем HTTP-код ответа
	assert.Equal(t, http.StatusOK, rec.Code)

	// Проверка содержимого ответа
	var actualInfo models.InfoResponse
	err := json.NewDecoder(rec.Body).Decode(&actualInfo)
	assert.NoError(t, err)
	assert.Equal(t, expectedInfo, &actualInfo) // Сравниваем структуру ответа

	// Убеждаемся, что все ожидания mock сработали
	mockRepo.AssertExpectations(t)
}

package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alias1177/merch-store/internal/config/config"
	"github.com/Alias1177/merch-store/internal/handlers/handlers"
	Jwtm "github.com/Alias1177/merch-store/internal/middleware/jwt" // исправлен импорт
	"github.com/Alias1177/merch-store/internal/models"
	"github.com/Alias1177/merch-store/internal/repositories"
	"github.com/Alias1177/merch-store/internal/usecase/auth"
	"github.com/Alias1177/merch-store/internal/usecase/buy"
	"github.com/Alias1177/merch-store/internal/usecase/coins"
	"github.com/Alias1177/merch-store/internal/usecase/info"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestServer(t *testing.T) *httptest.Server {
	cfg := config.Config{
		Database: config.DatabaseConfig{
			DSN: "host=localhost port=6000 user=myuser password=mypassword dbname=mydb sslmode=disable",
		},
		JWT: config.JWTConfig{
			Secret: "supersecretkey",
		},
	}

	ctx := context.Background()
	repo := repositories.New(ctx, cfg.Database.DSN)

	sendUsecase := coins.NewCoinsUsecase(repo)
	buyUsecase := buy.NewBuyUsecase(repo)
	infoUsecase := info.NewInfoUsecase(repo)
	userUsecase := auth.New(repo, cfg.JWT.Secret)

	handler := handlers.New(userUsecase, buyUsecase, infoUsecase, sendUsecase)

	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Post("/auth", handler.RegisterHandler)

		// Добавляем защищенные маршруты в отдельную группу
		r.Group(func(r chi.Router) {
			r.Use(Jwtm.JWTMiddleware(cfg.JWT.Secret))
			r.Get("/buy/{item}", handler.HandleBuy)
			r.Post("/sendCoin", handler.HandleSendCoins)
			r.Get("/info", handler.HandleInfo)
		})
	})

	return httptest.NewServer(r)
}

func TestPurchaseMerchFlow(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	// 1. Регистрация пользователя
	registerPayload := models.RegisterRequest{ // используем корректную структуру
		Username: "testuserr",
		Password: "testpasss",
	}
	jsonData, err := json.Marshal(registerPayload)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", ts.URL+"/api/auth", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)

	// Логируем тело ответа в случае ошибки
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Logf("Auth response: %s", string(body))
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var tokenResp models.TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	require.NoError(t, err)
	require.NotEmpty(t, tokenResp.Token)

	// 2. Покупка товара
	buyReq, err := http.NewRequest("GET", fmt.Sprintf("%s/api/buy/1", ts.URL), nil)
	require.NoError(t, err)
	buyReq.Header.Set("Authorization", "Bearer "+tokenResp.Token)

	resp, err = client.Do(buyReq)
	require.NoError(t, err)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Logf("Buy response: %s", string(body))
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSendCoinsFlow(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	client := &http.Client{}

	// 1. Регистрация отправителя
	senderPayload := models.RegisterRequest{ // используем корректную структуру
		Username: "sender",
		Password: "pass123",
	}
	jsonData, err := json.Marshal(senderPayload)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", ts.URL+"/api/auth", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var senderToken models.TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&senderToken)
	require.NoError(t, err)

	// 2. Регистрация получателя
	receiverPayload := models.RegisterRequest{ // используем корректную структуру
		Username: "receiver",
		Password: "pass123",
	}
	jsonData, err = json.Marshal(receiverPayload)
	require.NoError(t, err)

	req, err = http.NewRequest("POST", ts.URL+"/api/auth", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 3. Отправка монет
	sendCoinsPayload := models.SendCoinRequest{
		ToUser: "receiver",
		Amount: 100,
	}
	jsonData, err = json.Marshal(sendCoinsPayload)
	require.NoError(t, err)

	req, err = http.NewRequest("POST", ts.URL+"/api/sendCoin", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+senderToken.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	require.NoError(t, err)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Logf("SendCoin response: %s", string(body))
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

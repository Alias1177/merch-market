package main

import (
	"github.com/Alias1177/merch-store/internal/config/config"
	"github.com/Alias1177/merch-store/internal/handlers/handlers"
	Jwtm "github.com/Alias1177/merch-store/internal/middleware/jwt"
	"github.com/Alias1177/merch-store/internal/usecase/auth"
	"github.com/Alias1177/merch-store/internal/usecase/buy"
	"github.com/Alias1177/merch-store/internal/usecase/info"
	"log"
	"net/http"

	"github.com/Alias1177/merch-store/internal/repositories"
	"github.com/Alias1177/merch-store/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	logger.ColorLogger()
	cfg := config.Load(".env")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Подключение к БД
	repo := repositories.New(cfg.Database.DSN)
	defer repo.Close()

	buyUsecase := buy.NewBuyUsecase(repo)
	infoUsecase := info.NewInfoUsecase(repo)
	userUsecase := auth.New(repo, cfg.JWT.Secret)

	handler := handlers.New(userUsecase, buyUsecase, infoUsecase)

	r.Route("/api", func(route chi.Router) {
		route.Post("/auth", handler.RegisterHandler)

		route.Group(func(protected chi.Router) {
			protected.Use(Jwtm.JWTMiddleware(cfg.JWT.Secret))
			protected.Get("/buy/{item}", handler.HandleBuy)
			protected.Get("/info", handler.HandleInfo) // Добавляем endpoint
		})
	})

	log.Printf("Server running on port %s", cfg.App.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.App.Port, r))
}

//1)регистрация
//curl -X POST http://localhost:8080/api/auth \
//-H "Content-Type: application/json" \
//-d '{"username": "testuser", "password": "testpass"}'
//
//2)Покупка предмета по его ID
//curl -X GET http://localhost:8080/api/buy/1 \
//-H "Authorization: Bearer YOUR_JWT_TOKEN"
//
//3)отправка монет пользователю по его имени
//curl -X POST http://localhost:8080/api/sendCoin \
//-H "Content-Type: application/json" \
//-H "Authorization: Bearer YOUR_JWT_TOKEN" \
//-d '{"toUser": "anotheruser", "amount": 100}'
//
//4)получение информации о пользователе
//curl -X GET http://localhost:8080/api/info \
//-H "Authorization: Bearer YOUR_JWT_TOKEN"

//МИГРАЦИИ
//==применить все миграции==
//docker run --rm -v $(pwd)/migrations:/migrations migrate/migrate \
//  -path=/migrations \
//  -database "postgres://admin:secret@host.docker.internal:6000/mydatabase?sslmode=disable" up
//==откатить последнюю миграцию==
//docker run --rm -v $(pwd)/migrations:/migrations migrate/migrate \
//  -path=/migrations \
//  -database "postgres://admin:secret@host.docker.internal:6000/mydatabase?sslmode=disable" down 1
//  ==посмотреть версию текущей миграции==
//docker run --rm -v $(pwd)/migrations:/migrations migrate/migrate \
//  -path=/migrations \
//  -database "postgres://admin:secret@host.docker.internal:6000/mydatabase?sslmode=disable" version

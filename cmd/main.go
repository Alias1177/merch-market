package main

import (
	Jwtm "github.com/Alias1177/merch-store/internal/middleware/jwt"
	"log"
	"net/http"

	"github.com/Alias1177/merch-store/config/config"
	"github.com/Alias1177/merch-store/internal/handlers"
	"github.com/Alias1177/merch-store/internal/repositories"
	"github.com/Alias1177/merch-store/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	logger.ColorLogger()
	cfg := config.Load("./config/config.yaml")

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	// Подключение к БД
	db := repositories.Connect(cfg.Database.DSN)
	defer db.Close()

	r.Route("/api", func(route chi.Router) {
		// Публичные маршруты (без JWT)
		route.Post("/auth", handlers.RegisterHandler(db))

		// Защищенные маршруты (с JWT)
		route.Group(func(protected chi.Router) {
			protected.Use(Jwtm.JWTMiddleware(cfg.JWT.Secret))
			protected.Get("/buy/{item}", handlers.BuyHandler(db))
			protected.Post("/sendCoin", handlers.SendCoinsHandler(db))
			protected.Get("/info", handlers.InfoHandler(db))
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

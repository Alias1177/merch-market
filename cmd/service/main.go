package main

import (
	"github.com/Alias1177/merch-store/internal/config/config"
	"github.com/Alias1177/merch-store/internal/handlers/handlers"
	Jwtm "github.com/Alias1177/merch-store/internal/middleware/jwt"
	"github.com/Alias1177/merch-store/internal/usecase/auth"
	"github.com/Alias1177/merch-store/internal/usecase/buy"
	"github.com/Alias1177/merch-store/internal/usecase/coins"
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
	
	sendUsecase := coins.NewCoinsUsecase(repo)
	buyUsecase := buy.NewBuyUsecase(repo)
	infoUsecase := info.NewInfoUsecase(repo)
	userUsecase := auth.New(repo, cfg.JWT.Secret)

	handler := handlers.New(userUsecase, buyUsecase, infoUsecase, sendUsecase)

	r.Route("/api", func(route chi.Router) {
		route.Post("/auth", handler.RegisterHandler)

		route.Group(func(protected chi.Router) {
			protected.Use(Jwtm.JWTMiddleware(cfg.JWT.Secret))
			protected.Get("/buy/{item}", handler.HandleBuy)
			protected.Get("/info", handler.HandleInfo)
			protected.Post("/sendCoin", handler.HandleSendCoins)
		})
	})

	log.Printf("Server running on port %s", cfg.App.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.App.Port, r))
}

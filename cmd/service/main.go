package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alias1177/merch-store/internal/config/config"
	"github.com/Alias1177/merch-store/internal/handlers/handlers"
	Jwtm "github.com/Alias1177/merch-store/internal/middleware/jwt"
	"github.com/Alias1177/merch-store/internal/repositories"
	"github.com/Alias1177/merch-store/internal/usecase/auth"
	"github.com/Alias1177/merch-store/internal/usecase/buy"
	"github.com/Alias1177/merch-store/internal/usecase/coins"
	"github.com/Alias1177/merch-store/internal/usecase/info"
	"github.com/Alias1177/merch-store/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	logger.ColorLogger()

	ctx := context.Background()

	cfg := config.Load(os.Getenv("CONFIG_PATH"))

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Подключение к БД
	repo := repositories.New(ctx, cfg.Database.DSN)

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

	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second, // можно вынести в конфиг
		WriteTimeout: 10 * time.Second, // можно вынести в конфиг
		IdleTimeout:  30 * time.Second, // можно вынести в конфиг
	}

	// Запуск сервера в отдельной горутине
	go func() {
		slog.Info("starting server", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Ожидание сигнала для graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop
	slog.Info("received shutdown signal", "signal", sig.String())

	// Graceful shutdown с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server shutdown failed", "error", err)
	}

	slog.Info("server gracefully stopped")
}

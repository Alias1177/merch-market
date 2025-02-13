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

	cfg := config.Load(".env")

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
		Addr:    ":" + cfg.App.Port,
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	serverErrors := make(chan error, 1)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	select {
	case err := <-serverErrors:
		slog.Error("Ошибка запуска сервера", "error", err)

		if err := repo.Close(); err != nil {
			slog.Error("Ошибка при закрытии соединения с БД", "error", err)
		}
		os.Exit(1)
	case sig := <-quit:
		slog.Info("Получен сигнал завершения", "signal", sig)

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("Ошибка при graceful shutdown HTTP сервера", "error", err)
		}

		if err := repo.Close(); err != nil {
			slog.Error("Ошибка при закрытии соединения с БД", "error", err)
		}

		if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
			slog.Warn("Превышено время ожидания graceful shutdown")
		}
	}

	slog.Info("Сервер успешно остановлен")
}

package main

import (
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/slg"
	"url-shortener/internal/storage/sqlite"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofiber/fiber/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))
	log.Info("initializing server", slog.String("address", cfg.Address))
	log.Debug("logger debug mode is enabled")
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to initialize storage", slg.Err(err))
	}
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

}
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	}
	return log
}

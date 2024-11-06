package main

import (
	"log/slog"
	"net/http"
	"os"
	"server/internal/config"
	"server/internal/lib/logger/slogf"
	"server/internal/storage/sqlite"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
)

func main() {
	cfg := config.MustLoad() //НУ инициализация конфигурации проекта

	log := setupLogger(cfg.Env) // Ну вроде создание логгера

	log = log.With(slog.String("env", cfg.Env)) //Заставляет всегда в конце лога писать на какой именно машине было произведено выполнение проекта
	log.Info("starting url-shortner")
	log.Debug("debug messages are enabled")

	//Создание экземпляра, через который будем работать с бд. Но сейчас просто создали новую таблицу
	storage, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		log.Error("failed to init storage", slogf.Err(err))
		os.Exit(1)
	}
	_ = storage
	// storage.SaveURL("google.com", "sdffff")
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// router.Post("/url", save.New(log, storage))
	// router.Get("/urls", take.New(log, storage))
	// router.Get("/url/{id}", take.NewByID(log, storage))
	log.Info("starting server", slog.String("address", cfg.Address))

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
	log.Error("Error!")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}
	return log
}

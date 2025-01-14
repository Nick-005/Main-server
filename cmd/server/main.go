package main

import (
	"log/slog"
	"net/http"
	"os"
	"server/internal/config"
	"server/internal/lib/logger/slogf"
	"server/internal/server/handlers/auth"
	"server/internal/server/handlers/save"
	"server/internal/server/handlers/take"
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

	//Создание экземпляра, через который будем работать с бд. Но сейчас просто создали новую таблицу вакансий
	storageUser, err := sqlite.CreateTableUser(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", slogf.Err(err))
		os.Exit(1)
	}

	// storageUser, err := sqlite.CreateEmployeeTable()

	storageVac, err := sqlite.CreateVacancyTable(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", slogf.Err(err))
		os.Exit(1)
	}
	// Создание таблицы Работадателя
	storageEmp, err := sqlite.CreateEmployeeTable(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", slogf.Err(err))
		os.Exit(1)
	}
	// Ниже идёт создание роутера и использование его для middleware из chi
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/user", auth.NewUser(log, storageUser))       // POST запрос для добавления нового пользователя
	router.Post("/user/auth", auth.AuthUser(log, storageUser)) // POST запрос для авторизации пользователя по хэшу пароля + логина

	router.Post("/vac", save.NewVac(log, storageVac)) // POST запрос для добавления новой вакансии
	router.Post("/emp", save.NewEmp(log, storageEmp)) // POST запрос для добавления новой организации

	router.Get("/vac/{id}", take.GetVacancyByID(log, storageVac))  // GET запрос для получения данных о вакансии по её ID
	router.Get("/emp/{id}", take.GetEmployeeByID(log, storageEmp)) // GET запрос для получения данных о работадателе по его ID

	router.Get("/vacs", take.GetAllVacancy(log, storageVac))   // GET запрос для получения данных обо всех вакансиях
	router.Get("/emps", take.GetAllEmployees(log, storageEmp)) // GET запрос для получения данных обо всех работадателях
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

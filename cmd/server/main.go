package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/denissscare/todo-go/internal/config"
	gettodos "github.com/denissscare/todo-go/internal/handlers/getTodos"
	savetodo "github.com/denissscare/todo-go/internal/handlers/saveTodo"
	sqlite "github.com/denissscare/todo-go/internal/storage"
	"github.com/go-chi/chi/v5"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	var config = config.LoadConfig()
	var logger = setupLogger(config.Env)

	_ = logger

	fmt.Printf("Запуск сервера на: %s\n", config.Address)

	storage, err := sqlite.New(config.StoragePath)
	if err != nil {
		fmt.Printf("Ошибка инициализации БД: %s\n", err)
		os.Exit(1)
	}

	_ = storage

	router := chi.NewRouter()

	router.Post("/add-todo", savetodo.New(storage))
	router.Get("/all-todos", gettodos.New(storage))

	server := &http.Server{
		Addr:         config.Address,
		Handler:      router,
		ReadTimeout:  config.HTTPServer.Timeout,
		WriteTimeout: config.HTTPServer.Timeout,
		IdleTimeout:  config.HTTPServer.IdleTimeout,
	}
	fmt.Printf("Cервер запущен")
	if err := server.ListenAndServe(); err != nil {
		//TODO: Добавить логирование ERROR
		fmt.Printf("Ошибка запуска сервера: %s", err)
	}

	//TODO: Добавить логирование INFO
	fmt.Print("Сервер остановлен")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return logger
}

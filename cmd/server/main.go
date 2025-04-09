package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/denissscare/todo-go/internal/config"
	savetodo "github.com/denissscare/todo-go/internal/handlers/saveTodo"
	sqlite "github.com/denissscare/todo-go/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	var config = config.LoadConfig()

	fmt.Printf("Запуск сервера на: %s\n", config.Address)

	storage, err := sqlite.New(config.StoragePath)
	if err != nil {
		fmt.Printf("Ошибка инициализации БД: %s\n", err)
		os.Exit(1)
	}

	_ = storage

	router := chi.NewRouter()

	router.Post("/add-todo", savetodo.New(storage))

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

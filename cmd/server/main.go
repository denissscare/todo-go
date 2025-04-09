package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/denissscare/todo-go/internal/config"
	sqlite "github.com/denissscare/todo-go/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	var config = config.LoadConfig()

	fmt.Printf("Сервер запущен на: %s\n", config.Address)

	storage, err := sqlite.New(config.StoragePath)
	if err != nil {
		fmt.Println("Ошибка инициализации БД")
		os.Exit(1)
	}

	_ = storage

	router := chi.NewRouter()

	server := &http.Server{
		Addr:         config.Address,
		Handler:      router,
		ReadTimeout:  config.HTTPServer.Timeout,
		WriteTimeout: config.HTTPServer.Timeout,
		IdleTimeout:  config.HTTPServer.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		//TODO: Добавить логирование ERROR
		fmt.Printf("Ошибка запуска сервера: %s", err)
	}

	//TODO: Добавить логирование INFO
	fmt.Print("Сервер остановлен")
}

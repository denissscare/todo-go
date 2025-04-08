package main

import (
	"fmt"
	"net/http"

	"github.com/denissscare/todo-go/internal/config"
	"github.com/go-chi/chi/v5"
)

// TODO: Добавить логирование INFO
func startInfo(cfg *config.Config) {
	fmt.Printf("\nСервер запущен по адресу: %v\n", cfg.HTTPServer.Address)

	fmt.Print("\nПараметры конфига:")
	fmt.Printf("Storage path: %v\n", cfg.StoragePath)
	fmt.Printf("Address: %v\n", cfg.HTTPServer.Address)
	fmt.Printf("Timeout: %v\n\n\n", cfg.HTTPServer.Timeout)
}

func main() {
	var config = config.MustLoad()

	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("INFO: GET /\n")
		w.Write([]byte("Start page"))
	})

	server := &http.Server{
		Addr:         config.Address,
		Handler:      router,
		ReadTimeout:  config.HTTPServer.Timeout,
		WriteTimeout: config.HTTPServer.Timeout,
		IdleTimeout:  config.HTTPServer.IdleTimeout,
	}

	startInfo(config)
	if err := server.ListenAndServe(); err != nil {
		//TODO: Добавить логирование ERROR
		fmt.Printf("Ошибка запуска сервера: %s", err)
	}

	//TODO: Добавить логирование INFO
	fmt.Print("Сервер остановлен")
}

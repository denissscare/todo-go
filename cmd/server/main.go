package main

import (
	"fmt"

	"github.com/denissscare/todo-go/internal/config"
)

// TODO: Добавить логирование INFO
func printConfig(cfg *config.Config) {
	fmt.Printf("\n\nStorage path: %v\n", cfg.StoragePath)
	fmt.Printf("Address: %v\n", cfg.HTTPServer.Address)
	fmt.Printf("Timeout: %v\n\n\n", cfg.HTTPServer.Timeout)
}

func main() {
	var config = config.MustLoad()
	printConfig(config)
}

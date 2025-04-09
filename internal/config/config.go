package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:9090"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func getConfigPath() string {
	_, filename, _, _ := runtime.Caller(0)
	var rootDir string = filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	return filepath.Join(rootDir, "config", "env.yaml")
}

func LoadConfig() *Config {
	var configPath string = getConfigPath()

	//TODO: Добавить логирование ERROR.
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("Файл не найден: %s", configPath)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		fmt.Printf("Невозможно прочитать конфиг: %s", err)
	}

	return &config
}

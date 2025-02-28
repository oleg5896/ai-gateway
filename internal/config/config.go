package config

import (
	"github.com/joho/godotenv"
	"github.com/oleg5896/ai-common/logger"
	"os"
)

type Config struct {
	TelegramToken string
	GRPCInputAddr string
}

func Load() Config {
	err := godotenv.Load() // Загружаем .env файл, если есть
	if err != nil {
		logger.Info("Не удалось загрузить .env файл, используются переменные окружения")
	}

	return Config{
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		GRPCInputAddr: os.Getenv("GRPC_INPUT_ADDR"),
	}
}

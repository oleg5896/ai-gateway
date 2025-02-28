package main

import (
	"github.com/oleg5896/ai-common/logger"
	"github.com/oleg5896/ai-gateway/internal/config"
	"github.com/oleg5896/ai-gateway/internal/grpcclient"
	"github.com/oleg5896/ai-gateway/internal/handlers"
)

func main() {
	cfg := config.Load()
	if cfg.TelegramToken == "" || cfg.GRPCInputAddr == "" {
		logger.Error("Не заданы TELEGRAM_TOKEN или GRPC_INPUT_ADDR", nil)
		return
	}

	client, conn, err := grpcclient.NewInputClient(cfg.GRPCInputAddr)
	if err != nil {
		return
	}
	defer conn.Close()

	telegramHandler, err := handlers.NewTelegramHandler(cfg.TelegramToken, client)
	if err != nil {
		logger.Error("Ошибка создания Telegram-бота", err)
		return
	}

	telegramHandler.Start()
}

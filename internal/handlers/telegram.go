package handlers

import (
	"context"
	"github.com/oleg5896/ai-common/logger"
	proto "github.com/oleg5896/ai-proto/gateway"
	"time"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type TelegramHandler struct {
	bot    *telego.Bot
	client proto.GatewayClient
}

func NewTelegramHandler(token string, client proto.GatewayClient) (*TelegramHandler, error) {
	// Создаём бота с токеном и включаем отладку
	bot, err := telego.NewBot(token, telego.WithDefaultDebugLogger())
	if err != nil {
		return nil, err
	}

	return &TelegramHandler{
		bot:    bot,
		client: client,
	}, nil
}

func (h *TelegramHandler) Start() {
	logger.Info("Запуск Telegram-бота")

	// Настраиваем Long Polling
	params := &telego.GetUpdatesParams{
		Timeout: 60, // Таймаут в секундах
	}
	updates, err := h.bot.UpdatesViaLongPolling(params)
	if err != nil {
		logger.Error("Ошибка запуска Long Polling", err)
		return
	}
	defer h.bot.StopLongPolling()

	// Обрабатываем обновления
	for update := range updates {
		if update.Message == nil {
			continue
		}

		userName := update.Message.From.Username
		if userName == "" {
			userName = update.Message.From.FirstName // Если username отсутствует
		}
		logger.Info("Получено сообщение от " + userName + ": " + update.Message.Text)

		// Отправляем запрос в Input Processor через gRPC
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		req := &proto.Request{
			UserId: userName,
			Text:   update.Message.Text,
		}
		resp, err := h.client.ProcessRequest(ctx, req)
		if err != nil {
			logger.Error("Ошибка gRPC", err)
			continue
		}

		// Отправляем ответ пользователю
		chatID := telego.ChatID{ID: update.Message.Chat.ID}
		msg := tu.Message(chatID, "Статус: "+resp.Status)
		if _, err := h.bot.SendMessage(msg); err != nil {
			logger.Error("Ошибка отправки сообщения", err)
		}
	}
}

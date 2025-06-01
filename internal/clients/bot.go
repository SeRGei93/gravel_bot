package clients

import (
	"gravel_bot/internal/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func InitBot(cfg config.Bot) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		panic(err)
	}
	bot.Debug = true
	return bot
}

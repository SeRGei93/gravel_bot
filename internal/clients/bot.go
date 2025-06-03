package clients

import (
	"gravel_bot/internal/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var AwaitingMessage = make(map[int64]bool)

func InitBot(cfg config.Bot) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		panic(err)
	}
	SetLocalCommands(bot, cfg.AdminChat)
	bot.Debug = true
	return bot
}

func SetLocalCommands(bot *tgbotapi.BotAPI, chatID int64) error {
	commands := []tgbotapi.BotCommand{
		//{Command: "export", Description: "Список участников"},
		{Command: "export_csv", Description: "Список участников"},
		{Command: "export_gifts", Description: "Список подарков"},
		{Command: "send_notify", Description: "Отправить сообщение всем участникам"},
	}

	scope := tgbotapi.NewBotCommandScopeChat(chatID)

	cfg := tgbotapi.NewSetMyCommandsWithScopeAndLanguage(scope, "ru", commands...)

	_, err := bot.Request(cfg)
	return err
}

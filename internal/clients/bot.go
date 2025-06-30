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
	SetLocalCommands(bot, cfg)
	bot.Debug = true
	return bot
}

func SetLocalCommands(bot *tgbotapi.BotAPI, cfg config.Bot) error {
	adminCommands := []tgbotapi.BotCommand{
		{Command: "export_csv", Description: "Список участников"},
		{Command: "export_gifts", Description: "Список подарков"},
		{Command: "send_notify", Description: "Отправить сообщение всем участникам"},
	}

	err := sendCommandRequest(bot, cfg.AdminChat, adminCommands)
	if err != nil {
		return err
	}

	publicCommands := []tgbotapi.BotCommand{
		{Command: "public_info", Description: "Инфо"},
	}

	err = sendCommandRequest(bot, cfg.PublicChat, publicCommands)
	if err != nil {
		return err
	}

	return nil
}

func sendCommandRequest(bot *tgbotapi.BotAPI, chatId int64, commands []tgbotapi.BotCommand) error {
	scope := tgbotapi.NewBotCommandScopeChat(chatId)
	_, err := bot.Request(tgbotapi.NewSetMyCommandsWithScopeAndLanguage(scope, "ru", commands...))
	return err
}

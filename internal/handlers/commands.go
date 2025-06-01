package handlers

import (
	"gravel_bot/internal/database"
	"gravel_bot/internal/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Commands(bot *tgbotapi.BotAPI, update tgbotapi.Update, db database.Database) {
	switch update.Message.Command() {
	case "start":
		services.Start(bot, update, db)
	case "kamni200":
		services.Kamni200(bot, update, db)
	case "kamni200_off":
		services.Kamni200Off(bot, update, db)
	}
}

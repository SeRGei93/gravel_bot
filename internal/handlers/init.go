package handlers

import (
	"gravel_bot/internal/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Init(bot *tgbotapi.BotAPI, db database.Database) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	// Loop through each update.
	for update := range updates {
		if update.Message.IsCommand() {
			Commands(bot, update, db)
		}
	}
}

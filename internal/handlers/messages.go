package handlers

import (
	"gravel_bot/internal/await"
	"gravel_bot/internal/config"
	"gravel_bot/internal/database"
	"gravel_bot/internal/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Messages(bot *tgbotapi.BotAPI, update tgbotapi.Update, db database.Database, cfg config.Bot) {

	awaiting, exist := await.GetAwaiting(update.Message.From.ID)
	if !exist {
		return
	}

	switch awaiting.Type {
	case await.AwaitGift:
		services.SaveGift(bot, update, db, cfg)
	case await.AwaitResult:
		services.SaveResult(bot, update, db, cfg)
	default:
		services.NoHandler(bot, update, db, cfg)
	}
}

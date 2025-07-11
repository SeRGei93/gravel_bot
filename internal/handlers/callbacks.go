package handlers

import (
	"gravel_bot/internal/config"
	"gravel_bot/internal/database"
	"gravel_bot/internal/services"
	"gravel_bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Callbacks(bot *tgbotapi.BotAPI, update tgbotapi.Update, db database.Database, cfg config.Bot) {
	services.Start(bot, update, db, cfg)
	return

	cmd := update.CallbackQuery.Data
	switch cmd {
	case "rules":
		services.Rules(bot, update, db, cfg)
	case "kamni200":
		services.SetBike(bot, update, db, cfg)
	case "kamni200_off":
		services.Kamni200Off(bot, update, db, cfg)
	case "add_gift":
		services.AddGift(bot, update, db, cfg)
	case "add_result":
		services.AddResult(bot, update, db, cfg)
	default:
		_, bike := utils.GetKeyValue(update.CallbackQuery.Data)
		services.Kamni200(bot, update, db, cfg, bike)
	}
}

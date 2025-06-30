package services

import (
	"gravel_bot/internal/config"
	"gravel_bot/internal/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func PublicInfo(bot *tgbotapi.BotAPI, update tgbotapi.Update, db database.Database, cfg config.Bot) {
	if update.Message.Chat.ID == cfg.PublicChat {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Информация для участников")
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("🤖 Gravel Бот", "https://t.me/kamnigravelride_bot")),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("🏆 Призовой фонд", "https://t.me/kamnigravel/7698")),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("‼️ Условия участия", "https://t.me/kamnigravel/7697")),
		)

		bot.Send(msg)
	}
}

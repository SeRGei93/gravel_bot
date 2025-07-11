package services

import (
	"gravel_bot/internal/config"
	"gravel_bot/internal/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Start(bot *tgbotapi.BotAPI, update tgbotapi.Update, db database.Database, cfg config.Bot) {
	if update.Message.Chat.ID == cfg.PublicChat {
		return
	}

	text := `
<b>КАМНИ 200 🔥 18+</b>
Мероприятие завершилось!`

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	/*
		buttons, err := addButtons(update.Message.From.ID, "kamni200", db, cfg)
		if err == nil {
			msg.ReplyMarkup = buttons
		}*/
	msg.ParseMode = "HTML"

	bot.Send(msg)
}

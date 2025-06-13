package services

import (
	"gravel_bot/internal/await"
	"gravel_bot/internal/config"
	"gravel_bot/internal/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AddResult(bot *tgbotapi.BotAPI, update tgbotapi.Update, db database.Database, cfg config.Bot) {
	// пометить пользователя как ожидающего ввода
	await.SetAwaiting(update.CallbackQuery.From.ID, 600, await.AwaitResult)

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, `
Пожалуйста, отправьте ссылку на свой заезд следующим сообщением.

<b>Принимаются только ссылки на:</b>
• strava.com  
• komoot.com
`)
	msg.ParseMode = "HTML"
	bot.Send(msg)
}

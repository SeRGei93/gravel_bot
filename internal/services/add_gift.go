package services

import (
	"gravel_bot/internal/config"
	"gravel_bot/internal/database"
	"gravel_bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AddGift(bot *tgbotapi.BotAPI, update tgbotapi.Update, db database.Database, cfg config.Bot) {
	// пометить пользователя как ожидающего ввода
	utils.SetAwaiting(update.CallbackQuery.From.ID, 600)

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, `
	✏️ Укажите номинацию и опишите приз.

	Например:
	Первое место Топ кэп "спаси и сохрани"
	Книга цитат Стэтхэма за 8 место в абсолютном зачете
	За самый высокий средний пульс на дистанции упаковка мельдония
	Человек с самой лысой резиной получит блин шу пуэра

	❗Обязательно уложиться в одно сообщение
	`)
	msg.ParseMode = "HTML"
	bot.Send(msg)
}

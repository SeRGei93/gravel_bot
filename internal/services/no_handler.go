package services

import (
	"gravel_bot/internal/config"
	"gravel_bot/internal/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NoHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update, db database.Database, cfg config.Bot) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID

	msgText := `<b>Бот не ведёт диалог сообщениями.</b>
Все действия выполняются через кнопки.

Если хотите отправить результат заезда — нажмите кнопку
<b>«🏁 Я уже проехал»</b> и следуйте инструкциям.`

	msg := tgbotapi.NewMessage(chatID, msgText)
	msg.ParseMode = "HTML"

	// кнопки (если addButtons не вернёт ошибку)
	if buttons, err := addButtons(update.Message.From.ID, "kamni200", db, cfg); err == nil {
		msg.ReplyMarkup = buttons
	}

	_, _ = bot.Send(msg)
}

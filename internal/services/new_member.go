package services

import (
	"fmt"
	"gravel_bot/internal/config"
	"gravel_bot/internal/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewMember(bot *tgbotapi.BotAPI, update tgbotapi.Update, db database.Database, cfg config.Bot) {
	for _, newUser := range update.Message.NewChatMembers {
		if newUser.IsBot {
			continue // Игнорируем ботов
		}

		if update.Message.Chat.ID == cfg.PublicChat {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("👋 Привет, %s! Добро пожаловать в КАМНИ 200 🚴‍♂️", newUser.FirstName))
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("✅ Принять участие", "https://t.me/kamnigravelride_bot")),
				tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("🏆 Призовой фонд", "https://t.me/kamnigravel/7698")),
				tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("‼️ Условия участия", "https://t.me/kamnigravel/7697")),
			)

			bot.Send(msg)
		}
	}
}

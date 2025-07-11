package services

import (
	"fmt"
	"gravel_bot/internal/config"
	"gravel_bot/internal/database"
	"log"
	"regexp"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Dialog(bot *tgbotapi.BotAPI, update tgbotapi.Update, db database.Database, cfg config.Bot) {
	// Пропускаем не текстовые сообщения
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	userID := update.Message.From.ID
	chatId := update.Message.Chat.ID
	text := update.Message.Text
	adminChatId := cfg.AdminChat

	// Если сообщение от админа
	if chatId == adminChatId {
		if !strings.HasPrefix(text, "kamni=") {
			return
		}

		userID, msgText, err := handleAdminMessageRegex(text)
		if err != nil {
			msg := tgbotapi.NewMessage(adminChatId, err.Error())
			bot.Send(msg)
			return
		}

		msg := tgbotapi.NewMessage(userID, msgText)
		bot.Send(msg)

		return
	}

	// Если нет активного диалога - пересылаем сообщение админу
	msg := tgbotapi.NewMessage(adminChatId, fmt.Sprintf("kamni=%d", userID))
	bot.Send(msg)

	fwd := tgbotapi.NewForward(adminChatId, update.Message.Chat.ID, update.Message.MessageID)
	if _, err := bot.Send(fwd); err != nil {
		log.Printf("Ошибка пересылки сообщения: %v", err)
	}
}

func handleAdminMessageRegex(text string) (int64, string, error) {
	re := regexp.MustCompile(`^kamni=(\d+)\s+(.*)$`)
	matches := re.FindStringSubmatch(text)
	if len(matches) != 3 {
		return 0, "", fmt.Errorf("")
	}

	tgID, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, "", fmt.Errorf("некорректный ID пользователя: %v", err)
	}

	return tgID, matches[2], nil
}

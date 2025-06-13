package services

import (
	"fmt"
	"gravel_bot/internal/await"
	"gravel_bot/internal/config"
	"gravel_bot/internal/database"
	"log/slog"
	"regexp"
	"strings"

	"database/sql"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SaveResult(bot *tgbotapi.BotAPI, update tgbotapi.Update, db database.Database, cfg config.Bot) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	link := update.Message.Text

	validateResult := validateResultLink(link)
	if !validateResult.Valid {
		msg := tgbotapi.NewMessage(chatID, buildInvalidURLMessage(validateResult.Platform))
		msg.ParseMode = "HTML"
		bot.Send(msg)
		return
	}

	event, err := db.Event.FindEventByName("kamni200")
	if err != nil {
		slog.Error("ошибка: " + err.Error())
		bot.Send(tgbotapi.NewMessage(chatID, "Произошла ошибка"))
		return
	}

	userEvent, err := db.UserEvent.FindUserToEvent(userID, event.ID)
	if err != nil {
		slog.Error("ошибка: " + err.Error())
		bot.Send(tgbotapi.NewMessage(chatID, "Произошла ошибка"))
		return
	}

	userEvent.ResultLink = sql.NullString{String: link, Valid: true}

	err = db.UserEvent.UpdateUserEvent(userEvent)
	if err != nil {
		slog.Error("ошибка: " + err.Error())
		bot.Send(tgbotapi.NewMessage(chatID, "Произошла ошибка"))
		return
	}

	msg := tgbotapi.NewMessage(chatID, `Результат принят. Спасибо за участие 🫶`)
	msg.ParseMode = "HTML"
	buttons, err := addButtons(update.Message.From.ID, "kamni200", db, cfg)
	if err == nil {
		msg.ReplyMarkup = buttons
	}
	bot.Send(msg)

	await.DeleteAwaiting(userID)

	// переслать сообщение в админский чат
	notice := tgbotapi.NewMessage(cfg.AdminChat, fmt.Sprintf("@%s проехал КАМНИ 200 🏁", update.Message.From.UserName))
	bot.Send(notice)
	fwd := tgbotapi.NewForward(cfg.AdminChat, chatID, update.Message.MessageID)
	bot.Send(fwd)
}

type ValidationResult struct {
	Valid    bool
	Platform string
}

type ResultExample struct {
	Platform   string
	ExampleURL string
}

func validateResultLink(link string) ValidationResult {
	link = strings.TrimSpace(strings.ToLower(link))

	stravaRe := regexp.MustCompile(`^https?://(www\.)?strava\.com/activities/\d+$`)
	komootRe := regexp.MustCompile(`^https?://(www\.)?komoot\.com/tour/\d+$`)

	switch {
	case stravaRe.MatchString(link):
		return ValidationResult{Valid: true, Platform: "strava"}
	case komootRe.MatchString(link):
		return ValidationResult{Valid: true, Platform: "komoot"}
	case strings.Contains(link, "strava.com"):
		return ValidationResult{
			Valid:    false,
			Platform: "strava",
		}
	case strings.Contains(link, "komoot.com"):
		return ValidationResult{
			Valid:    false,
			Platform: "komoot",
		}
	default:
		return ValidationResult{
			Valid:    false,
			Platform: "",
		}
	}
}

func buildInvalidURLMessage(platform string) string {
	examples := []ResultExample{
		{Platform: "strava", ExampleURL: "https://www.strava.com/activities/14758223172"},
		{Platform: "komoot", ExampleURL: "https://www.komoot.com/tour/2308024419"},
	}

	caser := cases.Title(language.Und)

	text := "Ссылка не распознана.\n\nДопустимы только ссылки с платформ:\n"
	for _, result := range examples {
		text += fmt.Sprintf("- %s\n", caser.String(result.Platform))
	}

	for _, result := range examples {
		if result.Platform == platform {
			text += fmt.Sprintf("\nПример для %s: <code>%s</code>", caser.String(result.Platform), result.ExampleURL)
			return text
		}
	}

	// Если нет платформы — показать первый
	if len(examples) > 0 {
		result := examples[0]
		text += fmt.Sprintf("\nПример: <code>%s</code>", result.ExampleURL)
	}

	return text
}

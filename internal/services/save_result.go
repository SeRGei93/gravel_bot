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

	text := "Результат принят. Спасибо за участие 🫶"
	if validateResult.Platform == Strava {
		text += `

Не забудьте убирать ограничения приватности по стартовой и финишной точке в strava`
	}

	msg := tgbotapi.NewMessage(chatID, text)
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

type Platform string

const (
	Strava Platform = "strava"
	Komoot Platform = "komoot"
	None   Platform = ""
)

type ValidationResult struct {
	Valid    bool
	Platform Platform
}

type ResultExample struct {
	Platform   Platform
	ExampleURL string
}

func validateResultLink(link string) ValidationResult {
	link = strings.TrimSpace(strings.ToLower(link))

	stravaRe := regexp.MustCompile(`^https?://(www\.)?strava\.com/activities/\d+$`)
	stravaAppRe := regexp.MustCompile(`^https?://(www\.)?strava\.app\.link/[A-Za-z0-9]+$`)
	komootRe := regexp.MustCompile(`^https?://(www\.)?komoot\.com/tour/\d+$`)

	switch {
	case stravaRe.MatchString(link):
		return ValidationResult{Valid: true, Platform: Strava}
	case stravaAppRe.MatchString(link):
		return ValidationResult{Valid: true, Platform: Strava}
	case komootRe.MatchString(link):
		return ValidationResult{Valid: true, Platform: Komoot}
	case strings.Contains(link, "strava.com"):
		return ValidationResult{
			Valid:    false,
			Platform: Strava,
		}
	case strings.Contains(link, "komoot.com"):
		return ValidationResult{
			Valid:    false,
			Platform: Komoot,
		}
	default:
		return ValidationResult{
			Valid:    false,
			Platform: None,
		}
	}
}

func buildInvalidURLMessage(platform Platform) string {
	examples := []ResultExample{
		{Platform: Strava, ExampleURL: "https://www.strava.com/activities/14758223172"},
		{Platform: Komoot, ExampleURL: "https://www.komoot.com/tour/2308024419"},
	}

	caser := cases.Title(language.Und)

	text := "Ссылка не распознана.\n\nДопустимы только ссылки с платформ:\n"
	for _, result := range examples {
		text += fmt.Sprintf("- %s\n", caser.String(string(result.Platform)))
	}

	for _, result := range examples {
		if result.Platform == platform {
			text += fmt.Sprintf("\nПример для %s: <code>%s</code>", caser.String(string(result.Platform)), result.ExampleURL)
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

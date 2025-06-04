package services

import (
	"encoding/csv"
	"fmt"
	"gravel_bot/internal/config"
	"gravel_bot/internal/database"
	"gravel_bot/internal/database/repository"
	"gravel_bot/internal/utils"
	"log/slog"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ExportGifts(bot *tgbotapi.BotAPI, update tgbotapi.Update, db database.Database, cfg config.Bot) {
	// Найти событие
	event, err := db.Event.FindEventByName("kamni200")
	if err != nil {
		slog.Error("ошибка поиска события: " + err.Error())
		return
	}

	// Сформировать путь к временном файлу
	tmpFilePath := fmt.Sprintf("kamni200_gifts_%d.csv", time.Now().Unix())
	defer os.Remove(tmpFilePath)

	// Сгенерировать CSV
	rows, err := db.Gift.ExportGifts(event.ID)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при получении списка подарков")
		bot.Send(msg)
		return
	}

	err = makeFile(bot, rows, tmpFilePath)
	if err != nil {
		slog.Error("ошибка при экспорте подарков: " + err.Error())
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при экспорте подарков")
		bot.Send(msg)
		return
	}

	// Открыть файл
	file, err := os.Open(tmpFilePath)
	if err != nil {
		slog.Error("ошибка открытия файла: " + err.Error())
		return
	}
	defer file.Close()

	// Отправить файл как документ в чат админов
	fileReader := tgbotapi.FileReader{
		Name:   tmpFilePath,
		Reader: file,
	}
	doc := tgbotapi.NewDocument(cfg.AdminChat, fileReader)
	doc.Caption = "Подарки"

	if _, err := bot.Send(doc); err != nil {
		slog.Error("ошибка отправки файла: " + err.Error())
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось отправить файл")
		bot.Send(msg)
		return
	}
}

func makeFile(bot *tgbotapi.BotAPI, rows []repository.GiftDto, tmpFilePath string) error {
	file, err := os.Create(tmpFilePath)
	if err != nil {
		slog.Error("ошибка создания файла: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"tg_id", "username", "first_name", "last_name", "gift", "images", "date"})

	for _, row := range rows {
		files := ""
		for i, f := range row.Files {
			link, _ := utils.GetFileURL(bot, f.ID)
			if i > 0 {
				files += ", "
			}
			files += link
		}

		writer.Write([]string{
			fmt.Sprintf("%d", row.ID),
			"@" + row.NickName,
			row.FirstName,
			row.LastName,
			row.Content,
			files,
			row.CreatedAt.Format("02.01.2006 15:04:05"),
		})
	}

	return nil
}

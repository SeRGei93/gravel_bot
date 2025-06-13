package services

import (
	"fmt"
	"gravel_bot/internal/await"
	"gravel_bot/internal/config"
	"gravel_bot/internal/database"
	"gravel_bot/internal/database/table"
	"log/slog"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SaveGift(bot *tgbotapi.BotAPI, update tgbotapi.Update, db database.Database, cfg config.Bot) {
	if update.Message != nil && update.Message.Chat.IsPrivate() {
		userID := update.Message.From.ID
		hasMediaGroup := false

		awaiting, exist := await.GetAwaiting(userID)
		if exist {
			awaiting.Time = time.Now().Add(time.Duration(3) * time.Second) // меняем время ожидания, мы уже дождались ответ
			event, err := db.Event.FindEventByName("kamni200")
			if err != nil {
				slog.Error("ошибка поиска события: " + err.Error())
				return
			}

			var files []table.File
			if update.Message.Photo != nil && len(update.Message.Photo) > 0 {
				photo := update.Message.Photo[len(update.Message.Photo)-1]
				files = append(files, table.File{
					ID:   photo.FileID,
					Type: "photo",
				})
			}

			if update.Message.MediaGroupID != "" && len(files) > 0 {
				hasMediaGroup = true
				existGift, err := db.Gift.FindGiftByMediaGroup(update.Message.MediaGroupID)
				if err == nil {
					photo := files[0]
					file := table.File{
						ID:       photo.ID,
						Type:     "photo",
						EntityId: existGift.ID,
					}

					if err := db.File.CreateFile(file); err != nil {
						slog.Error(err.Error())
					}
					return
				}
			}

			// переслать сообщение в админский чат
			notice := tgbotapi.NewMessage(cfg.AdminChat, fmt.Sprintf("@%s добавил 🎁", update.Message.From.UserName))
			bot.Send(notice)
			fwd := tgbotapi.NewForward(cfg.AdminChat, update.Message.Chat.ID, update.Message.MessageID)
			bot.Send(fwd)

			text := update.Message.Text
			if len(files) > 0 {
				text = update.Message.Caption
			}

			// Создание или обновление пользователя
			_ = db.User.CreateUser(table.User{
				ID:        userID,
				NickName:  update.Message.From.UserName,
				FirstName: update.Message.From.FirstName,
				LastName:  update.Message.From.LastName,
			})

			// сохраняем подарок
			gift := table.Gift{
				UserID:       userID,
				EventID:      event.ID,
				Content:      text,
				MediaGroupId: update.Message.MediaGroupID,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				Files:        files,
			}

			if err := db.Gift.CreateGift(gift); err != nil {
				slog.Error(err.Error())
			}

			if !hasMediaGroup {
				await.DeleteAwaiting(userID) // очистить
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "✅ Спасибо, Ваше сообщение отправлено.")
			buttons, err := addButtons(update.Message.From.ID, "kamni200", db, cfg)
			if err == nil {
				msg.ReplyMarkup = buttons
			}

			bot.Send(msg)
		} else {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❗ Чтобы добавить еще один приз, нажмите кнопку «Добавить приз»."))
		}
	}
}

package services

import (
	"fmt"
	"gravel_bot/internal/config"
	"gravel_bot/internal/database"
	"gravel_bot/internal/database/table"
	"log/slog"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Rules(bot *tgbotapi.BotAPI, update tgbotapi.Update, db database.Database, cfg config.Bot) {
	const text = `<b>УСЛОВИЯ УЧАСТИЯ (ДИСКЛЕЙМЕР) ‼️</b>

Для обычного человека гравийная поездка на 200–300 км не является лёгкой прогулкой и требует хорошей физической и моральной подготовки, планирования питания и питья, а также наличия всего необходимого для ремонта велосипеда, оказания медпомощи и эвакуации себя. 

Участие в КАМНЯХ означает полное принятие следующих условий:

<b>1. Участники</b>: 18+, регистрация, взнос в призовой фонд и хотя бы частичное прохождение маршрута (подтвержденное ссылкой на Strava). Приоритет — у тех, кто проехал полностью.

<b>2. Питание и питьё 🍼🍔</b>: сам обеспечиваешь себя. Употребляй 50–100 г углеводов в час, немного белка. Пей каждые 15–20 мин. Имей запас воды и еды, пополняй по ходу.

<b>3. Безопасность 🤌</b>: твоя ответственность. Обязателен шлем, рабочий велик, фонари, аптечка, деньги. Не стартуй, если болен. При перегреве — охлаждайся. При плохом самочувствии — сойти.

<b>4. ПДД и законы 🚗</b>: соблюдай ПДД и законы РБ. Полная ответственность — на тебе.

<b>5. Проблемы на маршруте 🧚‍♀️</b>: рассчитывай только на себя. Имей заряженный телефон. Экстренные номера: 103, 112, 102. Пиши в чат, если нужна помощь.

<b>6. Сход с дистанции ⛔️</b>: добираешься сам. Такси, друзья, попутки, ЖД — сам решаешь.

<b>7. Ремонт и техобслуживание 🚴‍♀️</b>: исправный велосипед, инструменты, камеры, латки, червяки. Проколы — дело обычное.

<b>8. Навигация 🏞</b>: сам ориентируешься, следуешь треку. Заряжай навигатор, разметки нет.

<b>9. Новичкам 🍬</b>: езжай в компании — весело, безопасно, проще.

<b>10. Риски и ответственность 🤌</b>: ты участвуешь на свой страх и риск. Организаторы ответственности не несут. Участвуя — ты подтверждаешь согласие со всеми условиями.`

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
	buttons, err := addButtons(update.CallbackQuery.Message, "kamni200", db, cfg)
	if err == nil {
		msg.ReplyMarkup = buttons
	}

	msg.ParseMode = "HTML"
	if _, err := bot.Send(msg); err != nil {
		slog.Error(err.Error())
	}
}

func SetBike(bot *tgbotapi.BotAPI, update tgbotapi.Update, db database.Database, cfg config.Bot) {
	from := update.CallbackQuery.Message.Chat
	userID := from.ID

	// Найти событие по имени
	event, err := db.Event.FindEventByName("kamni200")
	if err != nil {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Ошибка: событие не найдено")
		buttons, err := addButtons(update.CallbackQuery.Message, "kamni200", db, cfg)
		if err == nil {
			msg.ReplyMarkup = buttons
		}
		bot.Send(msg)
		return
	}

	application, _ := db.UserEvent.FindUserToEvent(userID, event.ID)
	if application != nil {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Вы уже зарегистрированы")
		buttons, err := addButtons(update.CallbackQuery.Message, "kamni200", db, cfg)
		if err == nil {
			msg.ReplyMarkup = buttons
		}
		bot.Send(msg)
		return
	}

	text := "Какой у вас велосипед"
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Гравийник", "type=gravel"),
			tgbotapi.NewInlineKeyboardButtonData("МТБ", "type=mtb"),
			tgbotapi.NewInlineKeyboardButtonData("Фикс", "type=fixedgear"),
			tgbotapi.NewInlineKeyboardButtonData("Шоссейник", "type=gay"),
		),
	)

	if _, err := bot.Send(msg); err != nil {
		slog.Error(err.Error())
	}
}

func Kamni200(bot *tgbotapi.BotAPI, update tgbotapi.Update, db database.Database, cfg config.Bot, bike string) {
	from := update.CallbackQuery.Message.Chat
	userID := from.ID

	// Создание или обновление пользователя
	_ = db.User.CreateUser(table.User{
		ID:        userID,
		NickName:  from.UserName,
		FirstName: from.FirstName,
		LastName:  from.LastName,
	})

	// Найти событие по имени
	event, err := db.Event.FindEventByName("kamni200")
	if err != nil {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Ошибка: событие не найдено")
		buttons, err := addButtons(update.CallbackQuery.Message, "kamni200", db, cfg)
		if err == nil {
			msg.ReplyMarkup = buttons
		}
		bot.Send(msg)
		return
	}

	application, _ := db.UserEvent.FindUserToEvent(userID, event.ID)
	if application != nil {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Вы уже зарегистрированы")
		buttons, err := addButtons(update.CallbackQuery.Message, "kamni200", db, cfg)
		if err == nil {
			msg.ReplyMarkup = buttons
		}
		bot.Send(msg)
		return
	}

	// Зарегистрировать пользователя на событие
	err = db.UserEvent.RegisterUserToEvent(userID, event.ID, true, bike)
	if err != nil {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Вы уже зарегистрированы")
		buttons, err := addButtons(update.CallbackQuery.Message, "kamni200", db, cfg)
		if err == nil {
			msg.ReplyMarkup = buttons
		}
		bot.Send(msg)
		return
	}

	// Успешное сообщение
	text := `<b>Спасибо! Ваша заявка принята 🔥</b>

Не забудьте добавить <b>приз</b> 🎁  
Это <u>обязательное условие</u> — каждый участник что-то ставит, и кто-то получает ваш <b>приз</b>.  
Приз не должен быть дорогим — главное, чтобы он был <i>от души</i> ❤️  
Если передумали участвовать — просто дайте знать, это ок.
`
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
	buttons, err := addButtons(update.CallbackQuery.Message, "kamni200", db, cfg)
	msg.ParseMode = "HTML"
	if err == nil {
		msg.ReplyMarkup = buttons
	}
	if _, err := bot.Send(msg); err != nil {
		slog.Error(err.Error())
	}

	notification := tgbotapi.NewMessage(cfg.AdminChat, fmt.Sprintf("🚴 Новый участник: %s (@%s) \nТип: %s",
		from.FirstName+" "+from.LastName,
		from.UserName,
		bike,
	))

	if _, err := bot.Send(notification); err != nil {
		slog.Error(err.Error())
	}
}

func Kamni200Off(bot *tgbotapi.BotAPI, update tgbotapi.Update, db database.Database, cfg config.Bot) {
	from := update.CallbackQuery.Message.Chat
	userID := from.ID

	err := db.User.DeleteUser(userID)
	if err != nil {
		slog.Error(err.Error())
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Пользователь не найден")
		buttons, err := addButtons(update.CallbackQuery.Message, "kamni200", db, cfg)
		if err == nil {
			msg.ReplyMarkup = buttons
		}
		bot.Send(msg)
		return
	}

	// Найти событие по имени
	event, err := db.Event.FindEventByName("kamni200")
	if err != nil {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Ошибка: событие не найдено")
		buttons, err := addButtons(update.CallbackQuery.Message, "kamni200", db, cfg)
		if err == nil {
			msg.ReplyMarkup = buttons
		}
		bot.Send(msg)
		return
	}

	err = db.UserEvent.UnRegisterUserToEvent(userID, event.ID)
	if err != nil {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Заявка не найдена")
		buttons, err := addButtons(update.CallbackQuery.Message, "kamni200", db, cfg)
		if err == nil {
			msg.ReplyMarkup = buttons
		}
		bot.Send(msg)
		return
	}

	// Успешное сообщение
	text := "Заявка отменена"
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
	buttons, err := addButtons(update.CallbackQuery.Message, "kamni200", db, cfg)
	if err == nil {
		msg.ReplyMarkup = buttons
	}
	if _, err := bot.Send(msg); err != nil {
		slog.Error(err.Error())
	}
}

func addButtons(message *tgbotapi.Message, eventName string, db database.Database, cfg config.Bot) (*tgbotapi.InlineKeyboardMarkup, error) {
	from := message.Chat
	userID := from.ID

	event, err := db.Event.FindEventByName(eventName)
	if err != nil {
		return nil, err
	}

	var applicationBtn []tgbotapi.InlineKeyboardButton
	application, _ := db.UserEvent.FindUserToEvent(userID, event.ID)
	if application == nil {
		applicationBtn = append(applicationBtn, tgbotapi.NewInlineKeyboardButtonData("✅ Принять участие", "kamni200"))
	} else {
		applicationBtn = append(applicationBtn, tgbotapi.NewInlineKeyboardButtonData("😢 Отказаться от участия", "kamni200_off"))
	}

	result := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(applicationBtn...),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("➕🎁 Добавить приз", "add_gift")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("‼️ Условия участия", "rules")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("🚴‍♀️ Чат участников", "http://t.me/kamnigravel")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("🏆 Призовой фонд", "https://t.me/kamnigravel/7698")),
	)

	return &result, nil
}

func ExportCsv(bot *tgbotapi.BotAPI, update tgbotapi.Update, db database.Database, cfg config.Bot) {
	// Найти событие
	event, err := db.Event.FindEventByName("kamni200")
	if err != nil {
		slog.Error("ошибка поиска события: " + err.Error())
		return
	}

	// Сформировать путь к временном файлу
	tmpFilePath := fmt.Sprintf("kamni200_%d_%d.csv", event.ID, time.Now().Unix())
	defer os.Remove(tmpFilePath)

	// Сгенерировать CSV
	err = db.UserEvent.ExportEventParticipantsCSV(event.ID, tmpFilePath)
	if err != nil {
		slog.Error("ошибка при экспорте CSV: " + err.Error())
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при экспорте CSV")
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
	doc.Caption = "Список писок"

	if _, err := bot.Send(doc); err != nil {
		slog.Error("ошибка отправки файла: " + err.Error())
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось отправить файл")
		bot.Send(msg)
		return
	}
}

func SendNotify(bot *tgbotapi.BotAPI, update tgbotapi.Update, db database.Database, cfg config.Bot) {
	// msg := tgbotapi.NewMessage(update.Message.Chat.ID, "команда в разработке")
	// if _, err := bot.Send(msg); err != nil {
	// 	slog.Error(err.Error())
	// }

	if update.Message.IsCommand() && update.Message.Command() == "send_notify" {
		text := strings.TrimSpace(update.Message.CommandArguments())

		if text == "" {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Укажи текст сообщения после команды."))
			return
		}

		SendBroadcast(bot, cfg.AdminUsers, text)
	}
}

func SendBroadcast(bot *tgbotapi.BotAPI, users []int64, text string) {
	for _, user := range users {
		msg := tgbotapi.NewMessage(user, text)
		msg.ParseMode = "HTML"

		_, err := bot.Send(msg)
		if err != nil {
			slog.Warn("ошибка отправки", "user_id", user, "error", err)
		}
		// опционально: пауза, чтобы не попасть под rate limit
		time.Sleep(50 * time.Millisecond)
	}
}

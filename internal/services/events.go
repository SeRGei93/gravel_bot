package services

import (
	"gravel_bot/internal/database"
	"gravel_bot/internal/database/table"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Start(bot *tgbotapi.BotAPI, update tgbotapi.Update, db database.Database) {
	text := `
<b>КАМНИ 200 🔥 18+</b>
16–30 июня

Гравийный маршрут 200 км в формате гонки/бревета/покатухи — кому что ближе. Индивидуальное прохождение на условии самообеспечения. В зачёт принимаются страва-треки с окном прохождения с 16 по 30 июня 2025 года включительно.

<b>РЕГИСТРАЦИЯ УЧАСТНИКОВ</b> откроется в ближайшее время.

<b>ПРИЗОВОЙ ФОНД</b> формируется самими участниками, претендовать на призы могут только те, кто сделал вклад. Ставить можно любые новые вещи на любое место либо условие. Например: ведро раков на 75 место, проездной на трамвай поломавшему велик, пачка минералки тому, кто потеряет сознание — дайте волю фантазии. Не обязательно ставить свою квартиру. Любой донат от участников приветствуется как дань уважения гравийному сообществу.

<b>ОБЯЗАТЕЛЬНОЕ СНАРЯЖЕНИЕ:</b> исправный велик, шлем, ремкомплект, питание, вода, навигация, аптечка, передний и задний свет.
Легенда маршрута: 70% неасфальтированная поверхность.

🗺 <b>МАРШРУТ:</b> <a href="https://ehai.club/kamni/Kamni200_2025_v1.gpx">GPX</a> | <a href="https://nakarte.me/#m=10/54.26482/27.30927&l=Y&nktl=JBZ7YVT6aBOO5xd2fESKEQ">Nakarte</a>
❗️До старта возможны изменения

‼️ <a href="https://t.me/kamnigravel/7697">УСЛОВИЯ УЧАСТИЯ</a>
🍓 <a href="https://t.me/kamnigravel/7698">ПРИЗОВОЙ ФОНД</a>
📣 <a href="http://t.me/kamnigravel">Чат для участников</a>

<b>Готов принять твою заявку на участие.</b>
	`

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "HTML"
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

func Kamni200(bot *tgbotapi.BotAPI, update tgbotapi.Update, db database.Database) {
	from := update.Message.From
	userID := from.ID

	// Создание или обновление пользователя
	_ = db.User.CreateUser(table.User{
		Id:        userID,
		NickName:  from.UserName,
		FirstName: from.FirstName,
		LastName:  from.LastName,
	})

	// Найти событие по имени
	event, err := db.Event.FindEventByName("kamni200")
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка: событие не найдено")
		bot.Send(msg)
		return
	}

	// Зарегистрировать пользователя на событие
	err = db.UserEvent.RegisterUserToEvent(userID, event.ID, true)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы уже зарегистрированы")
		bot.Send(msg)
		return
	}

	// Успешное сообщение
	text := "Заявка принята"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	if _, err := bot.Send(msg); err != nil {
		slog.Error(err.Error())
	}
}

func Kamni200Off(bot *tgbotapi.BotAPI, update tgbotapi.Update, db database.Database) {
	from := update.Message.From
	userID := from.ID

	err := db.User.DeleteUser(userID)
	if err != nil {
		slog.Error(err.Error())
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пользователь не найден")
		bot.Send(msg)
		return
	}

	// Найти событие по имени
	event, err := db.Event.FindEventByName("kamni200")
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка: событие не найдено")
		bot.Send(msg)
		return
	}

	err = db.UserEvent.UnRegisterUserToEvent(userID, event.ID)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Заявка не найдена")
		bot.Send(msg)
		return
	}

	// Успешное сообщение
	text := "Заявка отменена"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	if _, err := bot.Send(msg); err != nil {
		slog.Error(err.Error())
	}
}

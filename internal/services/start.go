package services

import (
	"gravel_bot/internal/config"
	"gravel_bot/internal/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Start(bot *tgbotapi.BotAPI, update tgbotapi.Update, db database.Database, cfg config.Bot) {
	if update.Message.Chat.ID == cfg.PublicChat {
		return
	}

	text := `
<b>КАМНИ 200 🔥 18+</b>
13 июня по 7 июля

Гравийный маршрут 200 км в формате гонки/бревета/покатухи — кому что ближе. Индивидуальное прохождение на условии самообеспечения. В зачёт принимаются страва-треки с окном прохождения с 13 июня по 7 июля 2025 года включительно.

<b>ПРИЗОВОЙ ФОНД</b> формируется самими участниками, претендовать на призы могут только те, кто сделал вклад. Ставить можно любые новые вещи на любое место либо условие. Например: ведро раков на 75 место, проездной на трамвай поломавшему велик, пачка минералки тому, кто потеряет сознание — дайте волю фантазии. Не обязательно ставить свою квартиру. Любой донат от участников приветствуется как дань уважения гравийному сообществу.

<b>ОБЯЗАТЕЛЬНОЕ СНАРЯЖЕНИЕ:</b> исправный велик, шлем, ремкомплект, питание, вода, навигация, аптечка, передний и задний свет.
Легенда маршрута: 70% неасфальтированная поверхность.

🗺 <b>МАРШРУТ:</b> <a href="https://ehai.club/kamni/Kamni200_2025_v2.gpx">GPX</a> | <a href="https://nakarte.me/#m=14/54.16301/27.16962&l=O&nktl=Aloijnu-Zj1O93z84U743g">Nakarte</a>
	`

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	buttons, err := addButtons(update.Message.From.ID, "kamni200", db, cfg)
	if err == nil {
		msg.ReplyMarkup = buttons
	}
	msg.ParseMode = "HTML"

	bot.Send(msg)
}

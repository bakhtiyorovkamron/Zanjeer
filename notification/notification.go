package notification

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Send(message string) {
	bot, err := tgbotapi.NewBotAPI(("6119114967:AAF_s6_sQkPCk4qJMbpZSV81UqhYtFuZRBc"))
	if err != nil {
		return
	}
	msg := tgbotapi.NewMessage(5709226930, message)

	bot.Send(msg)
}

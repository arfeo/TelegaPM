package main

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

/**
 *
 *	Configure and send a new response message to the user
 *
 */

func SendResponseMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, output Output) {
	userID := message.From.ID

	// Create a new message config object
	m := tgbotapi.NewMessage(message.Chat.ID, output.Response)

	// Append inline keyboard (if applicable)
	if output.Buttons != nil {
		markup := tgbotapi.NewInlineKeyboardMarkup()

		for i := range output.Buttons {
			markup.InlineKeyboard = append(
				markup.InlineKeyboard,
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.InlineKeyboardButton {
						Text: 	output.Buttons[i].Text,
						URL: 	&output.Buttons[i].URL,
					},
				),
			)
		}

		m.ReplyMarkup = markup
	} else {

		// Append reply keyboard
		if !current[userID].Continuous {

			// If it's not a continuous input, append default keyboard
			m.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("/new"),
					tgbotapi.NewKeyboardButton("/view"),
					tgbotapi.NewKeyboardButton("/remove"),
				),
			)
		} else {

			// If it's a continuous input, append cancel keyboard
			m.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("/cancel"),
				),
			)
		}
	}

	// Additional output message set up
	m.ParseMode = "HTML"
	m.DisableWebPagePreview = true

	// Send message to Bot
	if _, err := bot.Send(m); err != nil {
		log.Println(err)
	}
}

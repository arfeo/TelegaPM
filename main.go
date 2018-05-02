package main

import (
	"log"
	"gopkg.in/telegram-bot-api.v4"
	"encoding/json"
)

func main() {

	// Get Bot configuration
	c, ok := ReadFromFile("bot.config")
	if ok {
		if len(c) > 0 {
			err := json.Unmarshal([]byte(c), &config)

			if err != nil {
				log.Println(err)
			}
		}
	}

	// Create a new Bot instance
	bot, err := tgbotapi.NewBotAPI(config.Token)

	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Get Bot updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	// Make global maps
	current = make(map[int]Current)
	masters = make(map[int]string)

	for update := range updates {

		// User sent a message to Bot
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			go RootInputProcessing(bot, update.Message)
		}
	}
}

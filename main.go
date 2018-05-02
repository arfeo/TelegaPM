package main

import (
	"log"
	"gopkg.in/telegram-bot-api.v4"
	"encoding/json"
)

func main() {

	// Get Bot configuration
	if !FileExists("bot.config") {
		log.Fatalln("Fatal error: Cannot find configuration file")
	}

	c, ok := ReadFromFile("bot.config")
	if !ok {
		log.Fatalln("Fatal error: Cannot read configuration file")
	}

	if err := json.Unmarshal([]byte(c), &config); err != nil {
		log.Fatalln(err)
	}

	if config.Token == "" {
		log.Fatalln("Fatal error: No BOT_TOKEN found in the configuration file")
	}

	// Create a new Bot instance
	bot, err := tgbotapi.NewBotAPI(config.Token)

	if err != nil {
		log.Fatalln(err)
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

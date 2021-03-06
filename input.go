package main

import (
	"crypto/sha256"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"strconv"
)

/**
 * Function processes root level input
 */
func RootInputProcessing(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	var (
		userID        = message.From.ID
		storage       []Storage
		storageName   = strconv.Itoa(userID) + "." + config.Extension
		output        Output
		hash          = sha256.Sum256([]byte(masters[userID]))
	)

	// `whoami` command is always available
	if message.Text != "/whoami" {

		// If Bot lock is set up in config -- check the match
		if len(config.Lock) > 0 {
			if ok, _ := ValueInArray(strconv.Itoa(userID), config.Lock); !ok {
				BotSendMsg(bot, message, "This Bot is locked for private use only.")

				return
			}
		}

		// Check for storage existence & master password (if needed)
		if message.Text != "/start" && message.Text != "/cancel" && current[userID].Command != "/start" && !current[userID].Continuous {

			// If the storage does not exist, user should use `/start` command
			if !FileExists(storageName) {
				delete(masters, userID)

				BotSendMsg(bot, message, "Storage not found. Please, use /start command to set up Bot.")

				return
			}

			// If `masters[]` for the user is not defined, do `/start`
			if _, me := masters[userID]; !me {
				message.Text = "/start"
			}
		}
	}

	switch message.Text {

	// Initialize Bot
	case "/start" :
		{
			// Prompt for master password if masters[] for the user not set
			if _, ok := masters[userID]; ok {
				output.Response = "Hello again."

				StopContinuousInput(message)
			} else {
				output = StartContinuousInput(bot, message)
			}
			break
		}

	// Drop the storage
	case "/drop" :
		{
			BotSendMsg(bot, message, "Drop the storage...")

			output = StartContinuousInput(bot, message)
			break
		}

	// Create a new element
	case "/new" :
		{
			BotSendMsg(bot, message, "Create a new element...")

			output = StartContinuousInput(bot, message)
			break
		}

	// View an element
	case "/view" :
		{
			// Get elements from the storage file
			f, _ := ReadFromFile(storageName)
			storage = GetStorage(f, hash)

			b, ok := StorageList(storage)
			if ok {
				BotSendMsg(bot, message, "View an element...")

				// Print elements list from storage
				BotSendMsg(bot, message, b)

				output = StartContinuousInput(bot, message)
			} else {

				// Elements list is empty, return the respective message
				output.Response = b

				StopContinuousInput(message)
			}
			break
		}

	// Remove an element
	case "/remove" :
		{
			// Get elements from the storage file
			f, _ := ReadFromFile(storageName)
			storage = GetStorage(f, hash)

			b, ok := StorageList(storage)
			if ok {
				BotSendMsg(bot, message, "Remove an element...")

				// Print elements list from storage
				BotSendMsg(bot, message, b)

				output = StartContinuousInput(bot, message)
			} else {

				// Elements list is empty, return the respective message
				output.Response = b

				StopContinuousInput(message)
			}
			break
		}

	// Print the current user's ID
	case "/whoami" :
		{
			BotSendMsg(bot, message, strconv.Itoa(userID))
			break
		}

	// Cancel current input process
	case "/cancel" :
		{
			output.Response = "Canceled."

			StopContinuousInput(message)
			break
		}

	// Other
	default :
		{
			if current[userID].Continuous {
				output = ContinuousInputProcessing(bot, message)
			} else {
				output.Response = "Unknown command."

				StopContinuousInput(message)
			}
			break
		}
	}

	SendResponseMessage(bot, message, output)
}

/**
 * Function processes continuous input
 */
func ContinuousInputProcessing(bot *tgbotapi.BotAPI, message *tgbotapi.Message) Output {
	var (
		userID        = message.From.ID
		storage       []Storage
		storageName   = strconv.Itoa(userID) + "." + config.Extension
		output        Output
		hash          = sha256.Sum256([]byte(masters[userID]))
	)

	bc := current[userID]

	switch bc.Command {

	// Initialize Bot (continuous)
	case "/start" :
		{
			switch bc.Step {
			case 1 :
				{
					bc.Step++
					current[userID] = bc

					if !FileExists(storageName) {

						// If storage for the user does not exist, print welcome message
						inline := tgbotapi.NewMessage(message.Chat.ID,
							"<code>TelegaPM</code> is an open source Telegram password management bot engine.\n\n" +
							"WARNING: <code>TelegaPM</code> is designed as a self-hosted service. " +
							"Since one's credentials is a sensitive data, please, consider this Bot as the engine demo, " +
							"and <strong>DO NOT STORE</strong> anything private in it.")
						inline.ParseMode = "HTML"

						if _, err := bot.Send(inline); err != nil {
							log.Println(err)
						}

						output.Response = "Please, enter the master password for the storage:"
					} else {
						output.Response = "Master password:"
					}
					break
				}
			case 2 :
				{
					// Store master password
					masters[userID] = message.Text

					// Set hash with the new master password
					hash = sha256.Sum256([]byte(masters[userID]))

					if !FileExists(storageName) {

						// If storage for the user does not exist, create it
						if CreateFile(storageName) {

							// Create control element
							control := Current {
								Entry: Storage {
									ID: strconv.Itoa(userID),
								},
							}

							// Save control element to the storage file
							j := PushToStorage(storage, control, hash)
							if WriteToFile(storageName, j) {
								output.Response = "Storage successfully created. Master password is set"
							} else {
								RemoveFile(storageName)
								delete(masters, userID)
								output.Response = "Storage could not be created this time. Try again"
							}
						} else {
							delete(masters, userID)
							output.Response = "Storage could not be created this time. Try again"
						}
					} else {

						// Get elements from the storage file
						f, _ := ReadFromFile(storageName)
						buffer := GetStorage(f, hash)

						if len(buffer) > 0 {
							if c, _ := strconv.Atoi(buffer[0].ID); c == userID {
								storage = buffer
								output.Response = "Master password accepted."
							} else {
								delete(masters, userID)
								output.Response = "Master password incorrect."
							}
						} else {
							delete(masters, userID)
							output.Response = "Master password incorrect."
						}
					}

					StopContinuousInput(message)
					break
				}
			}
			break
		}

	// Drop the storage (continuous)
	case "/drop" :
		{
			switch bc.Step {
			case 1 :
				{
					bc.Step++
					current[userID] = bc
					output.Response = "Are you sure you want to drop the storage? This action cannot be undone. [Y/n]"
					break
				}
			case 2 :
				{
					if message.Text == "Y" || message.Text == "y" {
						if RemoveFile(storageName) {
							output.Response = "Storage dropped."
						}
					} else {
						output.Response = "Canceled."
					}

					StopContinuousInput(message)
					break
				}
			}
			break
		}

	// Create a new element (continuous)
	case "/new" :
		{
			switch bc.Step {
			case 1 :
				{
					// Store ID
					if u, err := uuid.NewV4(); err == nil {
						bc.Entry.ID = u.String()
					}

					bc.Step++
					current[userID] = bc
					output.Response = "Title:"
					break
				}
			case 2 :
				{
					// Store title
					bc.Entry.Title = message.Text

					bc.Step++
					current[userID] = bc
					output.Response = "Login:"
					break
				}
			case 3 :
				{
					// Store login
					bc.Entry.Login = message.Text

					bc.Step++
					current[userID] = bc
					output.Response = "Password:"
					break
				}
			case 4 :
				{
					// Store password
					bc.Entry.Pass = message.Text

					bc.Step++
					current[userID] = bc
					output.Response = "Url:"
					break
				}
			case 5 :
				{
					// Store url
					bc.Entry.Url = message.Text

					// Get elements from the storage file
					f, _ := ReadFromFile(storageName)
					storage = GetStorage(f, hash)

					// Save element to the storage file
					j := PushToStorage(storage, bc, hash)
					if WriteToFile(storageName, j) {
						output.Response = "Element successfully saved."
					} else {
						output.Response = "Error occurred while saving changes to file. Try again."
					}

					StopContinuousInput(message)
					break
				}
			}
			break
		}

	// View an element (continuous)
	case "/view" :
		{
			switch bc.Step {
			case 1 :
				{
					bc.Step++
					current[userID] = bc
					output.Response = "Element #:"
					break
				}
			case 2 :
				{
					// Store element ID
					bc.Element, _ = strconv.Atoi(message.Text)

					// Get elements from the storage file
					f, _ := ReadFromFile(storageName)
					storage = GetStorage(f, hash)

					b, p, ok := GetStorageElementInfo(storage, bc.Element)
					if ok {

						// Print the element info
						inline := tgbotapi.NewMessage(message.Chat.ID, b)
						inline.ParseMode = "HTML"

						if _, err := bot.Send(inline); err != nil {
							log.Println(err)
						}

						// Print password in a separate message
						output.Response = "<code>" + p + "</code>"
					} else {

						// Element not found, return the respective message
						output.Response = b
					}

					StopContinuousInput(message)
					break
				}
			}
			break
		}

	// Remove an element (continuous)
	case "/remove" :
		{
			switch bc.Step {
			case 1 :
				{
					bc.Step++
					current[userID] = bc
					output.Response = "Element #:"
					break
				}
			case 2 :
				{
					// Store element ID
					bc.Element, _ = strconv.Atoi(message.Text)

					// Get elements from the storage file
					f, _ := ReadFromFile(storageName)
					storage = GetStorage(f, hash)

					j,ok := RemoveElementFromStorage(storage, bc.Element, hash)
					if ok {

						// Update storage file
						if WriteToFile(storageName, j) {
							output.Response = "Element successfully removed."
						} else {
							output.Response = "Error occurred while saving changes to file. Try again."
						}
					} else {
						output.Response = "Element not found."
					}

					StopContinuousInput(message)
					break
				}
			}
			break
		}
	}

	return output
}

/**
 * Function sets current fields to start continuous input
 */
func StartContinuousInput(bot *tgbotapi.BotAPI, message *tgbotapi.Message) Output {
	b := current[message.From.ID]
	b = Current{
		Command:      message.Text,
		Step:         1,
		Continuous:   true,
	}

	current[message.From.ID] = b
	output := ContinuousInputProcessing(bot, message)

	return output
}

/**
 * Function resets current fields to stop continuous input
 */
func StopContinuousInput(message *tgbotapi.Message) {
	b := current[message.From.ID]
	b = Current{
		Command:      "",
		Step:         0,
		Continuous:   false,
	}

	current[message.From.ID] = b
}

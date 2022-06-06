package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
)

type Bot struct {
	bot         *tgbotapi.BotAPI
	updatesChan tgbotapi.UpdatesChannel
}

var WatchersList map[int64]int64
var Robot *Bot
var Keyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Subscribe To Updates"),
	),
)

func init() {

	WatchersList = make(map[int64]int64)
	bot, err := tgbotapi.NewBotAPI(os.Getenv("STINGERBOT"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	Robot = &Bot{updatesChan: updates, bot: bot}

}
func (b *Bot) SendAlert(message string) {
	for _, ChatID := range WatchersList {
		msg := tgbotapi.NewMessage(ChatID, message)
		msg.ReplyMarkup = Keyboard
		b.bot.Send(msg)
	}

}

func (b *Bot) StartLifeCycle(wg *sync.WaitGroup) {
	defer wg.Done()
	for update := range b.updatesChan {
		if update.Message != nil { // If we got a message
			log.Info("[%s] %s", update.Message.From.UserName, update.Message.Text)
			if update.Message.Text == "Subscribe To Updates" {
				WatchersList[update.Message.From.ID] = update.Message.Chat.ID
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Updated chat where I will notify you")
				msg.ReplyMarkup = Keyboard
				b.bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You can Only Subscribe")
				msg.ReplyMarkup = Keyboard
				b.bot.Send(msg)
			}

		}
	}
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/telegram-bot-api.v4"
)

const (
	botToken = "paste your bot token here"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
		Hello World!
	`))
}

func updateHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	if update.Message.Text == "/help" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пришли мне фото (документом) с подписью вида \n"+
			"верхний текст\n"+
			"@\n"+
			"нижний текст\n")
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	} else if update.Message.Document == nil ||
		update.Message.Caption == "" ||
		!strings.Contains(update.Message.Caption, "@") {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Прикольно. Но не понятно."+update.Message.Caption)
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	} else {
		text := strings.Split(update.Message.Caption, "@")
		text[0] = strings.Replace(text[0], "\n", "", -1)
		text[1] = strings.Replace(text[1], "\n", "", -1)

		photo := update.Message.Document

		url, err := bot.GetFileDirectURL(photo.FileID)
		if err != nil {
			fmt.Println("Error while getting url " + err.Error())
		}

		path := addText(url, text)

		msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, path)

		_, err = bot.Send(msg)
		if err != nil {
			fmt.Println("Error while sending meme " + err.Error())
		}
	}
}
func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/", handler)

	fmt.Println("starting server at :" + port)
	go http.ListenAndServe(":"+port, nil)

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		go updateHandler(update, bot)
	}
}

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/telebot.v3"
)

func parseDate(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func calculateDaysLeft() int {
	layout := "2006-01-02T15:04:05.000Z"
	birthDate, err := time.Parse(layout, os.Getenv("BIRTH_DATE"))
	if err != nil {
		log.Fatal(err)
	}
	parsedBirthDate := parseDate(birthDate.Year(), int(birthDate.Month()), birthDate.Day())
	deathDate := parsedBirthDate.AddDate(80, 0, 0)
	parsedDeathDate := parseDate(deathDate.Year(), int(deathDate.Month()), deathDate.Day())
	diff := parsedDeathDate.Sub(parsedBirthDate).Hours() / 24
	return int(diff)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading environment variables")
	}

	chatId, err := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	token := os.Getenv("TELEGRAM_TOKEN")
	if err != nil {
		log.Fatal(err)
	}

	conf := telebot.Settings{
		Token: token,
		Poller: &telebot.LongPoller{
			Timeout: 10 * time.Second,
		},
	}

	bot, err := telebot.NewBot(conf)
	if err != nil {
		log.Fatal(err)
	}

	chat := &telebot.Chat{ID: chatId}
	daysLeft := calculateDaysLeft()
	fmt.Println(daysLeft)

	bot.Send(chat, fmt.Sprintf("You have *%d* days left", daysLeft), &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdownV2,
	})

	bot.Start()
}

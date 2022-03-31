package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	github "memento-mori/pkg/github"

	"gopkg.in/telebot.v3"
)

func parseDate(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func calculateDaysLeft() int {
	birthDate, err := time.Parse(time.RFC3339, os.Getenv("BIRTH_DATE"))
	if err != nil {
		log.Fatal(err)
	}
	parsedBirthDate := parseDate(birthDate.Year(), int(birthDate.Month()), birthDate.Day())
	deathDate := parsedBirthDate.AddDate(80, 0, 0)
	parsedDeathDate := parseDate(deathDate.Year(), int(deathDate.Month()), deathDate.Day())
	diff := parsedDeathDate.Sub(time.Now()).Hours() / 24
	return int(diff)
}

func main() {
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

	events, err := github.GetLastPushEvents(os.Getenv("USER"))
	if err != nil {
		log.Fatal(err)
	}

	// TODO: sending only number of push events for now, send more data later
	bot.Send(chat, fmt.Sprintf("☠️ You have *%d* days left\n👨🏻‍💻 You made *%d* pushes to github yesterday", daysLeft, len(events)), &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdownV2,
	})

}

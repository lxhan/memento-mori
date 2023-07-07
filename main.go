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

func calculateDaysLeft() int {
	birthDate, err := time.Parse("2006-01-02", os.Getenv("BIRTH_DATE"))
	if err != nil {
		log.Fatal(err)
	}

	deathDate := birthDate.AddDate(80, 0, 0)

	difference := time.Until(deathDate).Hours()

	days := int(difference / 24)

	return days
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
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

	events, err := GetLastPushEvents(os.Getenv("USER"))
	if err != nil {
		log.Fatal(err)
	}

	gevents, err := GetGcalEvents()
	if err != nil {
		log.Fatal(err)
	}

	var gcalEvents string
	for _, e := range gevents {
		gcalEvents += fmt.Sprintf("• %s at `%s`\n", e["summary"], e["date"])
	}

	bot.Send(
		chat,
		fmt.Sprintf(
			"☠️ You have *%d* days left\n👨🏻‍💻 You made *%d* pushes to github yesterday\n\n📆 Upcoming events:\n%v",
			daysLeft,
			len(events),
			gcalEvents,
		),
		&telebot.SendOptions{
			ParseMode: telebot.ModeMarkdownV2,
		},
	)
}

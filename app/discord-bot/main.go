package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/meir/discord-bot/internal/bot"
	"github.com/meir/discord-bot/internal/commands"
	"github.com/meir/discord-bot/internal/events"
	"github.com/meir/discord-bot/internal/logging"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			logging.Warn(err)
			return
		}
		logging.Println("Shutting down.")
	}()

	b := bot.NewDiscordBot()
	b.RegisterEvents()

	err := b.Session.Open()
	if err != nil {
		logging.Fatal(err)
	}
	defer b.Session.Close()

	commands.RegisterCommands(b.Session)
	events.Register(b.Session)
	defer commands.RemoveCommands(b.Session)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	logging.Println(fmt.Sprintf("Started version #%v", os.Getenv("VERSION")))
	<-stop
}

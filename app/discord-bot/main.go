package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/meir/discord-bot/internal/bot"
	"github.com/meir/discord-bot/internal/commands"
	"github.com/meir/discord-bot/internal/events"
)

func main() {
	b := bot.NewDiscordBot()
	b.RegisterEvents()

	err := b.Session.Open()
	if err != nil {
		panic(err)
	}
	defer b.Session.Close()

	commands.RegisterCommands(b.Session)
	events.Register(b.Session)
	defer commands.RemoveCommands(b.Session)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}

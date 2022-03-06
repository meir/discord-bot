package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/discord-bot/internal/bot"
	"github.com/meir/discord-bot/internal/commands"
)

func main() {
	b := bot.NewDiscordBot()
	b.RegisterEvents()

	err := b.Session.Open()
	if err != nil {
		panic(err)
	}
	defer b.Session.Close()

	b.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
		IdleSince:  nil,
		Activities: []*discordgo.Activity{},
		AFK:        false,
		Status:     fmt.Sprintf("Running version #%v", os.Getenv("VERSION")),
	})

	commands.RegisterCommands(b.Session)
	defer commands.RemoveCommands(b.Session)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}

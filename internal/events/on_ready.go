package events

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

func OnReady(session *discordgo.Session, ready *discordgo.Ready) {
	session.UpdateStatusComplex(discordgo.UpdateStatusData{
		IdleSince:  nil,
		Activities: []*discordgo.Activity{},
		AFK:        false,
		Status:     fmt.Sprintf("Running version #%v", os.Getenv("VERSION")),
	})
}

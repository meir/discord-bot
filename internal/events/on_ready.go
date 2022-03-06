package events

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

func OnReady(session *discordgo.Session, ready *discordgo.Ready) {
	session.UpdateGameStatus(0, fmt.Sprintf("version #%v", os.Getenv("VERSION")))
}

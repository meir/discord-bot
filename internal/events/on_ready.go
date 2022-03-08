package events

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/discord-bot/internal/logging"
)

func OnReady(session *discordgo.Session, ready *discordgo.Event) {
	err := session.UpdateGameStatus(0, fmt.Sprintf("Running version %v", os.Getenv("VERSION")))
	// activity := &discordgo.Activity{
	// 	Name: fmt.Sprintf("Running version %v", os.Getenv("VERSION")),
	// 	Type: discordgo.ActivityTypeCustom,
	// }

	// err := session.UpdateStatusComplex(discordgo.UpdateStatusData{
	// 	Activities: []*discordgo.Activity{activity},
	// 	Status:     string(discordgo.StatusOnline),
	// })
	if err != nil {
		logging.Warn("Failed to update status", err)
	}
}

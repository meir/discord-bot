package utils

import (
	"github.com/bwmarrin/discordgo"
)

func HiddenResponse(session *discordgo.Session, interaction *discordgo.InteractionCreate, text string) {
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   1 << 6,
			Content: text,
		},
	})
}

func NoResponse(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Command ran succesfully",
		},
	})
	session.InteractionResponseDelete(session.State.User.ID, interaction.Interaction)
}

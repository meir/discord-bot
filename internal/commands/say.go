package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/meir/discord-bot/internal/logging"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	appendCommand("Say", "say", "I'll say whatever you want :)", []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "message",
			Description: "Message for me to say",
			Required:    true,
		},
	}, say)
}

func say(session *discordgo.Session, interaction *discordgo.InteractionCreate, db *mongo.Database) {
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponsePong,
	})

	session.ChannelMessageSend(interaction.ChannelID, interaction.ApplicationCommandData().Options[0].StringValue())

	err := session.InteractionResponseDelete(session.State.User.ID, interaction.Interaction)
	if err != nil {
		logging.Warn(err)
	}
}

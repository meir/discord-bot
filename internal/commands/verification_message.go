package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/discord-bot/internal/logging"
	"github.com/meir/discord-bot/internal/utils"
	"github.com/meir/discord-bot/pkg/structs"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	appendCommand("Verification Message", "verification-message", "Used to send the verification message in this channel", []*discordgo.ApplicationCommandOption{}, verification_message)
}

func verification_message(session *discordgo.Session, interaction *discordgo.InteractionCreate, db *mongo.Database) {
	if !utils.IsModerator(session, interaction, db) {
		utils.HiddenResponse(session, interaction, fmt.Sprintf("You don't have permission to use this command!"))
		return
	}

	query := structs.NewQuery(session, interaction, db)
	guild, err := query.Guild(interaction.GuildID)
	if err == mongo.ErrNoDocuments {
		guild = query.NewGuild(interaction.GuildID)
	} else {
		message, err := session.ChannelMessage(guild.VerificationMessage.ChannelID, guild.VerificationMessage.MessageID)
		if err == nil {
			session.ChannelMessageDelete(message.ChannelID, message.ID)
		}
	}

	msg, err := session.ChannelMessageSendComplex(interaction.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:       "Server Verification",
			Description: `In order to get access to the server you can make a verification request using this message.`,
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Request Verification",
						Style:    discordgo.PrimaryButton,
						CustomID: "verification_request_button",
					},
				},
			},
		},
	})
	if err != nil {
		logging.Warn(err)
		return
	}

	channel, err := session.Channel(msg.ChannelID)
	if err != nil {
		logging.Warn(err)
		return
	}

	guild.VerificationMessage = structs.Message{
		MessageID: msg.ID,
		ChannelID: channel.ID,
		ParentID:  channel.ParentID,
	}

	err = guild.Update()
	if err != nil {
		logging.Warn(err)
		return
	}

	utils.HiddenResponse(session, interaction, "Verification message has been placed.")
}

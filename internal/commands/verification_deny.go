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
	appendCommand("Verification Deny", "verification-deny", "Use this command to deny the current users verification channel.", []*discordgo.ApplicationCommandOption{
		{
			Name:        "reason",
			Description: "Reason to deny the verification",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    false,
		},
	}, verification_deny)
}

func verification_deny(session *discordgo.Session, interaction *discordgo.InteractionCreate, db *mongo.Database) {
	if !utils.IsModerator(session, interaction, db) {
		utils.HiddenResponse(session, interaction, fmt.Sprintf("You don't have permission to use this command!"))
		return
	}

	query := structs.NewQuery(session, interaction, db)

	channel, err := query.Channel(interaction.ChannelID)
	if err != nil {
		logging.Warn(err)
		return
	}

	var userId string
	if channelType, ok := channel.Metadata["type"]; ok && channelType == string(structs.VERIFICATION_CHANNEL) {
		userId = channel.Metadata["user"]
	}

	channel.Delete()
	session.ChannelDelete(channel.ChannelID)

	var reason string
	if len(interaction.ApplicationCommandData().Options) > 0 {
		reason = interaction.ApplicationCommandData().Options[0].StringValue()
		goto WithReason
	}

	session.GuildMemberDelete(interaction.GuildID, userId)
	goto RemoveChannel
WithReason:
	session.GuildMemberDeleteWithReason(interaction.GuildID, userId, reason)

RemoveChannel:
	session.ChannelDelete(channel.ChannelID)
}

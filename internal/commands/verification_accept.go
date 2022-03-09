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
	appendCommand("Verification Accept", "verification-accept", "Use this command to accept the current users verification channel.", []*discordgo.ApplicationCommandOption{}, verification_accept)
}

func verification_accept(session *discordgo.Session, interaction *discordgo.InteractionCreate, db *mongo.Database) {
	if !utils.IsModerator(session, interaction, db) {
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   1 << 6,
				Content: fmt.Sprintf("You don't have permission to use this command!"),
			},
		})
		return
	}

	query := structs.NewQuery(session, interaction, db)

	guild, err := query.Guild(interaction.GuildID)
	if err != nil {
		logging.Warn(err)
		return
	}

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
	session.GuildMemberRoleAdd(interaction.GuildID, userId, guild.VerifiedRole)
	session.ChannelDelete(channel.ChannelID)
}

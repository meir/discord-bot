package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/discord-bot/internal/utils"
	"github.com/meir/discord-bot/pkg/structs"
	"go.mongodb.org/mongo-driver/bson"
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
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   1 << 6,
				Content: fmt.Sprintf("You don't have permission to use this command!"),
			},
		})
		return
	}

	channels := db.Collection(os.Getenv("COLLECTION_CHANNELS"))

	// channel := interaction.ChannelID
	channelDocument := channels.FindOne(context.Background(), bson.M{
		"guild_id":   interaction.GuildID,
		"channel_id": interaction.ChannelID,
	})

	var channel structs.Channel
	err := channelDocument.Decode(&channel)
	if err != nil {
		panic(err)
	}

	var userId string
	if channelType, ok := channel.Metadata["type"]; ok && channelType == string(structs.VERIFICATION_CHANNEL) {
		userId = channel.Metadata["user"]
	}

	session.ChannelDelete(channel.ChannelID)
	channels.DeleteOne(context.Background(), channel)

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

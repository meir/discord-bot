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

	channels := db.Collection(os.Getenv("COLLECTION_CHANNELS"))
	guilds := db.Collection(os.Getenv("COLLECTION_SERVERDATA"))

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

	guildDocument := guilds.FindOne(context.Background(), bson.M{
		"guild_id": interaction.GuildID,
	})

	var guild structs.Guild
	err = guildDocument.Decode(&guild)
	if err != nil {
		panic(err)
	}

	var userId string
	if channelType, ok := channel.Metadata["type"]; ok && channelType == string(structs.VERIFICATION_CHANNEL) {
		userId = channel.Metadata["user"]
	}

	channels.DeleteOne(context.Background(), channel)
	session.GuildMemberRoleAdd(interaction.GuildID, userId, guild.VerifiedRole)
	session.ChannelDelete(channel.ChannelID)
}

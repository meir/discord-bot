package commands

import (
	"context"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/discord-bot/pkg/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	appendCommand("Verification Message", "verification-message", "Used to send the verification message in this channel", []*discordgo.ApplicationCommandOption{}, verification_message)
}

func verification_message(session *discordgo.Session, interaction *discordgo.InteractionCreate, db *mongo.Database) {
	col := db.Collection(os.Getenv("COLLECTION_SERVERDATA"))

	guildDocument := col.FindOne(context.Background(), bson.M{
		"guild_id": interaction.GuildID,
	})

	var message *discordgo.Message
	var guild structs.Guild
	err := guildDocument.Decode(&guild)
	if err == mongo.ErrNoDocuments {
		goto sendMessage
	}
	if err != nil {
		panic(err)
	}

	message, err = session.ChannelMessage(guild.VerificationMessage.ChannelID, guild.VerificationMessage.MessageID)
	if err != nil {
		goto sendMessage
	}
	session.ChannelMessageDelete(message.ChannelID, message.ID)

sendMessage:
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
		panic(err)
	}

	channel, err := session.Channel(msg.ChannelID)
	if err != nil {
		panic(err)
	}

	opts := options.Update().SetUpsert(true)
	col.UpdateOne(context.Background(), bson.M{
		"guild_id": interaction.GuildID,
	}, bson.M{
		"$set": bson.M{
			"verification_message": bson.M{
				"channel_id": msg.ChannelID,
				"message_id": msg.ID,
				"parent_id":  channel.ParentID,
			},
		},
	}, opts)

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   1 << 6,
			Content: "Verification message has been placed.",
		},
	})
}

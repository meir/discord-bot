package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/discord-bot/pkg/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	appendComponent(discordgo.ButtonComponent, "verification_request_button", verification_request_button)
}

func verification_request_button(session *discordgo.Session, interaction *discordgo.InteractionCreate, db *mongo.Database) {
	// Check if person already has an open verification channel
	channels := db.Collection(os.Getenv("COLLECTION_CHANNELS"))
	col := db.Collection(os.Getenv("COLLECTION_SERVERDATA"))

	overwrites := []*discordgo.PermissionOverwrite{
		{
			ID:    interaction.Member.User.ID,
			Type:  discordgo.PermissionOverwriteTypeMember,
			Allow: discordgo.PermissionViewChannel,
		},
		{
			ID:    interaction.Member.User.ID,
			Type:  discordgo.PermissionOverwriteTypeMember,
			Allow: discordgo.PermissionSendMessages,
		},
		{
			ID:   interaction.GuildID, // > The @everyone role has the same ID as the guild it belongs to.
			Type: discordgo.PermissionOverwriteTypeRole,
			Deny: discordgo.PermissionViewChannel,
		},
		{
			ID:   interaction.GuildID, // ^^^
			Type: discordgo.PermissionOverwriteTypeRole,
			Deny: discordgo.PermissionSendMessages,
		},
	}

	opts := options.Update().SetUpsert(true)

	var guildDocument *mongo.SingleResult
	var guild structs.Guild

	channelDocument := channels.FindOne(context.Background(), bson.M{
		"guild_id": interaction.GuildID,
		"metadata": bson.M{
			"type": structs.VERIFICATION_CHANNEL,
			"user": interaction.Member.User.ID,
		},
	})

	var channel structs.Channel
	var ch *discordgo.Channel
	err := channelDocument.Decode(&channel)
	if err == mongo.ErrNoDocuments {
		goto CreateChannel
	}
	if err != nil {
		panic(err)
	}

	ch, err = session.Channel(channel.ChannelID)
	if err != nil {
		// panic(err)
		goto CreateChannel
	}

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   1 << 6,
			Content: fmt.Sprintf("You already have an open verification channel: %v", ch.Mention()),
		},
	})
	goto EditChannel
CreateChannel:
	// Create verification channel
	guildDocument = col.FindOne(context.Background(), bson.M{
		"guild_id": interaction.GuildID,
	})

	err = guildDocument.Decode(&guild)
	if err != nil {
		panic(err)
	}

	ch, err = session.GuildChannelCreateComplex(interaction.GuildID, discordgo.GuildChannelCreateData{
		Name:                 fmt.Sprintf("verification-%v", interaction.Member.User.ID),
		Type:                 discordgo.ChannelTypeGuildText,
		Topic:                "You can explain here why you should be allowed into the server! Good luck :)",
		PermissionOverwrites: overwrites,
		ParentID:             guild.VerificationMessage.ParentID,
	})

	channels.UpdateOne(context.Background(), bson.M{
		"guild_id": ch.GuildID,
		"metadata": bson.M{
			"type": string(structs.VERIFICATION_CHANNEL),
			"user": interaction.Member.User.ID,
		},
	}, bson.M{
		"$set": structs.Channel{
			GuildID:   ch.GuildID,
			ChannelID: ch.ID,
			Metadata: map[string]string{
				"type": string(structs.VERIFICATION_CHANNEL),
				"user": interaction.Member.User.ID,
			},
		},
	}, opts)

	//session.GuildChannelCreate(interaction.GuildID, fmt.Sprintf("verification-%v", interaction.Member.User.ID), discordgo.ChannelTypeGuildText)
	if err != nil {
		panic(err)
	}

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   1 << 6,
			Content: fmt.Sprintf("Your verification channel has been created: %v", ch.Mention()),
		},
	})

EditChannel:
	session.ChannelEditComplex(ch.ID, &discordgo.ChannelEdit{
		PermissionOverwrites: overwrites,
	})
}

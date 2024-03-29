package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/discord-bot/internal/logging"
	"github.com/meir/discord-bot/internal/utils"
	"github.com/meir/discord-bot/pkg/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	appendCommand("Quote Channel", "quote-channel", "Sets this channel as the quotes channel.", []*discordgo.ApplicationCommandOption{}, quote_channel)
}

func quote_channel(session *discordgo.Session, interaction *discordgo.InteractionCreate, db *mongo.Database) {
	if !utils.IsModerator(session, interaction, db) {
		utils.HiddenResponse(session, interaction, fmt.Sprintf("You don't have permission to use this command!"))
		return
	}

	query := structs.NewQuery(session, interaction, db)

	channel, err := query.ChannelByFilter(bson.M{
		"guild_id": interaction.GuildID,
		"metadata": bson.M{
			"type": string(structs.QUOTES_CHANNEL),
		},
	})
	if err == mongo.ErrNoDocuments {
		channel = query.NewChannel(interaction.GuildID, interaction.ChannelID)
		err = nil
	}
	if err != nil {
		logging.Warn(err)
		return
	}

	channel.Metadata = map[string]string{
		"type": string(structs.QUOTES_CHANNEL),
	}
	channel.ChannelID = interaction.ChannelID
	err = channel.Update()
	if err != nil {
		logging.Warn(err)
		return
	}

	utils.HiddenResponse(session, interaction, fmt.Sprintf("This channel has been saved as the quotes channel."))
}

package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/discord-bot/internal/logging"
	"github.com/meir/discord-bot/pkg/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	appendCommand("Quote", "quote", "Adds a quote in the quote channel", []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "quote",
			Description: "The quote of the person",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "user",
			Description: "Person who said the quote",
			Required:    true,
		},
	}, quote)
}

func quote(session *discordgo.Session, interaction *discordgo.InteractionCreate, db *mongo.Database) {
	query := structs.NewQuery(session, interaction, db)

	count, err := query.QuoteCount(interaction.GuildID)
	if err == mongo.ErrNoDocuments {
		count = 0
		err = nil
	}
	if err != nil {
		logging.Warn("couldnt find quote count", err)
		return
	}

	quote := query.NewQuote(interaction.GuildID, count)
	quote.Message = interaction.ApplicationCommandData().Options[0].StringValue()
	quote.User = interaction.ApplicationCommandData().Options[1].UserValue(session).ID
	if quote.Update() != nil {
		logging.Warn("couldnt update quote", err)
		return
	}

	channel, err := query.ChannelByFilter(bson.M{
		"guild_id": interaction.GuildID,
		"metadata": bson.M{
			"type": string(structs.QUOTES_CHANNEL),
		},
	})

	if err != nil {
		logging.Warn("couldnt find quote channel", err)
		return
	}
	session.ChannelMessageSend(channel.ChannelID, fmt.Sprintf("\"%v\" - <@%v>", quote.Message, quote.User))
}

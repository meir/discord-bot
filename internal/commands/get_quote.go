package commands

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/discord-bot/internal/logging"
	"github.com/meir/discord-bot/pkg/structs"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	appendCommand("Get Quote", "get-quote", "Gets a quote in the quote channel", []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "quote-id",
			Description: "Get a specific quote or leave empty for a random one.",
			Required:    false,
		},
	}, get_quote)
}

func get_quote(session *discordgo.Session, interaction *discordgo.InteractionCreate, db *mongo.Database) {
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

	rand.Seed(time.Now().UnixMicro())
	id := rand.Int63n(count)
	if len(interaction.ApplicationCommandData().Options) > 0 {
		id = interaction.ApplicationCommandData().Options[0].IntValue()
	}

	quote, err := query.Quote(interaction.GuildID, id)
	if err != nil {
		logging.Warn(err)
		return
	}

	session.ChannelMessageSend(interaction.ChannelID, fmt.Sprintf("\"%v\" - <@%v>", quote.Message, quote.User))
}

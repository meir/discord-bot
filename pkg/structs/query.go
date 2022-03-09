package structs

import (
	"context"
	"os"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Query struct {
	session     *discordgo.Session
	interaction *discordgo.InteractionCreate
	db          *mongo.Database
}

func NewQuery(session *discordgo.Session, interaction *discordgo.InteractionCreate, db *mongo.Database) *Query {
	return &Query{session, interaction, db}
}

func (q *Query) NewGuild(id string) *Guild {
	g := &Guild{GuildID: id}
	g.Query = q
	return g
}

func (q *Query) Guild(id string) (*Guild, error) {
	guilds := q.db.Collection(os.Getenv("COLLECTION_SERVERDATA"))
	guildDoc := guilds.FindOne(context.Background(), bson.M{
		"guild_id": id,
	})

	var guild Guild
	err := guildDoc.Decode(&guild)
	if err != nil {
		return nil, err
	}
	guild.Query = q
	return &guild, nil
}

func (q *Query) GuildByFilter(filter bson.M) (*Guild, error) {
	guilds := q.db.Collection(os.Getenv("COLLECTION_SERVERDATA"))
	guildDoc := guilds.FindOne(context.Background(), filter)

	var guild Guild
	err := guildDoc.Decode(&guild)
	if err != nil {
		return nil, err
	}
	guild.Query = q
	return &guild, nil
}

func (q *Query) NewChannel(guildId string, id string) *Channel {
	c := &Channel{GuildID: id}
	c.Query = q
	return c
}

func (q *Query) Channel(id string) (*Channel, error) {
	channels := q.db.Collection(os.Getenv("COLLECTION_CHANNELS"))
	channelDoc := channels.FindOne(context.Background(), bson.M{
		"channel_id": id,
	})

	var channel Channel
	err := channelDoc.Decode(&channel)
	if err != nil {
		return nil, err
	}
	channel.Query = q
	return &channel, nil
}

func (q *Query) ChannelByFilter(filter bson.M) (*Channel, error) {
	channels := q.db.Collection(os.Getenv("COLLECTION_CHANNELS"))
	channelDoc := channels.FindOne(context.Background(), filter)

	var channel Channel
	err := channelDoc.Decode(&channel)
	if err != nil {
		return nil, err
	}
	channel.Query = q
	return &channel, nil
}

func (q *Query) NewQuote(guildID string, id int64) *Quote {
	g := &Quote{Number: id, GuildID: guildID}
	g.Query = q
	return g
}

func (q *Query) Quote(guildID string, id int64) (*Quote, error) {
	quotes := q.db.Collection(os.Getenv("COLLECTION_QUOTES"))
	channelDoc := quotes.FindOne(context.Background(), bson.M{
		"number":   id,
		"guild_id": guildID,
	})

	var quote Quote
	err := channelDoc.Decode(&quote)
	if err != nil {
		return nil, err
	}
	quote.Query = q
	return &quote, nil
}

func (q *Query) QuoteByFilter(filter bson.M) (*Quote, error) {
	quotes := q.db.Collection(os.Getenv("COLLECTION_QUOTES"))
	channelDoc := quotes.FindOne(context.Background(), filter)

	var quote Quote
	err := channelDoc.Decode(&quote)
	if err != nil {
		return nil, err
	}
	quote.Query = q
	return &quote, nil
}

func (q *Query) QuoteCount(guildID string) (int64, error) {
	quotes := q.db.Collection(os.Getenv("COLLECTION_QUOTES"))

	return quotes.CountDocuments(context.Background(), bson.M{
		"guild_id": guildID,
	})
}

package structs

import (
	"context"
	"errors"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Guild struct {
	*Query `json:"_", bson:"_"`

	GuildID             string  `bson:"guild_id"`
	VerifiedRole        string  `bson:"verified_role"`
	ModRole             string  `bson:"moderator_role"`
	VerificationMessage Message `bson:"verification_message"`
}

func (g *Guild) Update() error {
	if g.Query == nil {
		return errors.New("Cannot update with a non-query guild struct")
	}
	if g.GuildID == "" {
		return errors.New("Cannot use Update without a GuildID")
	}

	guilds := g.Query.db.Collection(os.Getenv("COLLECTION_SERVERDATA"))
	opts := options.Update().SetUpsert(true)
	_, err := guilds.UpdateOne(context.Background(), bson.M{
		"guild_id": g.GuildID,
	}, bson.M{
		"$set": g,
	}, opts)
	return err
}

type Message struct {
	*Query `json:"_", bson:"_"`

	ChannelID string `bson:"channel_id"`
	MessageID string `bson:"message_id"`
	ParentID  string `bson:"parent_id"`
}

type Channel struct {
	*Query `json:"_", bson:"_"`

	GuildID   string            `bson:"guild_id"`
	ChannelID string            `bson:"channel_id"`
	Metadata  map[string]string `bson:"metadata"`
}

func (c *Channel) Update() error {
	if c.Query == nil {
		return errors.New("Cannot update with a non-query guild struct")
	}
	if c.GuildID == "" {
		return errors.New("Cannot use Update without a ChannelID")
	}

	channels := c.Query.db.Collection(os.Getenv("COLLECTION_CHANNELS"))
	opts := options.Update().SetUpsert(true)
	_, err := channels.UpdateOne(context.Background(), bson.M{
		"channel_id": c.ChannelID,
	}, bson.M{
		"$set": c,
	}, opts)
	return err
}

func (c *Channel) Delete() error {
	if c.Query == nil {
		return errors.New("Cannot delete with a non-query guild struct")
	}
	if c.GuildID == "" {
		return errors.New("Cannot use delete without a GuildID")
	}
	if c.ChannelID == "" {
		return errors.New("Cannot use delete without a ChannelID")
	}

	channels := c.Query.db.Collection(os.Getenv("COLLECTION_CHANNELS"))
	_, err := channels.DeleteOne(context.Background(), bson.M{
		"channel_id": c.ChannelID,
		"guild_id":   c.GuildID,
	})
	return err
}

type ChannelType string

const (
	VERIFICATION_CHANNEL ChannelType = "verification_channel"
	QUOTES_CHANNEL       ChannelType = "quotes_channel"
)

type Quote struct {
	*Query `json:"-", bson:"-"`

	Number  int64  `bson:"number"`
	Message string `bson:"message"`
	User    string `bson:"user_id"`
	GuildID string `bson:"guild_id"`
}

func (q *Quote) Update() error {
	if q.Query == nil {
		return errors.New("Cannot update with a non-query guild struct")
	}
	if q.GuildID == "" {
		return errors.New("Cannot use Update without a GuildID")
	}

	quotes := q.Query.db.Collection(os.Getenv("COLLECTION_QUOTES"))
	opts := options.Update().SetUpsert(true)
	_, err := quotes.UpdateOne(context.Background(), bson.M{
		"number": q.Number,
	}, bson.M{
		"$set": q,
	}, opts)
	return err
}

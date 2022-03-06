package utils

import (
	"context"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/discord-bot/pkg/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func IsModerator(session *discordgo.Session, interaction *discordgo.InteractionCreate, db *mongo.Database) bool {
	servers := db.Collection(os.Getenv("COLLECTION_SERVERDATA"))
	serverDocument := servers.FindOne(context.Background(), bson.M{
		"guild_id": interaction.GuildID,
	})

	var server structs.Guild
	err := serverDocument.Decode(&server)
	if err != nil {
		panic(err)
	}

	for _, v := range interaction.Member.Roles {
		if v == server.ModRole {
			return true
		}
	}
	return false
}

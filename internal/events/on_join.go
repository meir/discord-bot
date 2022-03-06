package events

import (
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
)

func Handler(database *mongo.Database) interface{} {
	return func(s *discordgo.Session, i *discordgo.GuildMemberAdd) {

	}
}

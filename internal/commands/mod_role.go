package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	appendCommand("Moderator Role", "mod-role", "Used to set the moderator role.", []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionRole,
			Name:        "mod-role",
			Description: "Moderator Role",
			Required:    true,
		},
	}, mod_role)
}

func mod_role(session *discordgo.Session, interaction *discordgo.InteractionCreate, db *mongo.Database) {
	col := db.Collection(os.Getenv("COLLECTION_SERVERDATA"))
	opts := options.Update().SetUpsert(true)

	role := interaction.ApplicationCommandData().Options[0].RoleValue(session, interaction.GuildID)

	col.UpdateOne(context.Background(), bson.M{
		"guild_id": interaction.GuildID,
	}, bson.M{
		"$set": bson.M{
			"moderator_role": role.ID,
		},
	}, opts)

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   1 << 6,
			Content: fmt.Sprintf("The moderator role is now set to %v.", role.Mention()),
		},
	})
}

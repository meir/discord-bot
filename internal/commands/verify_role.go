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
	appendCommand("Verify Role", "verify-role", "Used to set the verification role.", []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionRole,
			Name:        "verification-role",
			Description: "Verification Role",
			Required:    true,
		},
	}, verify_role)
}

func verify_role(session *discordgo.Session, interaction *discordgo.InteractionCreate, db *mongo.Database) {
	col := db.Collection(os.Getenv("COLLECTION_SERVERDATA"))
	opts := options.Update().SetUpsert(true)

	role := interaction.ApplicationCommandData().Options[0].RoleValue(session, interaction.GuildID)

	col.UpdateOne(context.Background(), bson.M{
		"guild_id": interaction.GuildID,
	}, bson.M{
		"$set": bson.M{
			"verified_role": role.ID,
		},
	}, opts)

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   1 << 6,
			Content: fmt.Sprintf("The verified role is now set to %v.", role.Mention()),
		},
	})
}

package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/discord-bot/internal/logging"
	"github.com/meir/discord-bot/internal/utils"
	"github.com/meir/discord-bot/pkg/structs"
	"go.mongodb.org/mongo-driver/mongo"
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
	if g, err := session.Guild(interaction.GuildID); err != nil || g.OwnerID != interaction.User.ID {
		utils.HiddenResponse(session, interaction, "You have no permission for this command")
		return
	}
	role := interaction.ApplicationCommandData().Options[0].RoleValue(session, interaction.GuildID)

	guild, err := structs.NewQuery(session, interaction, db).Guild(interaction.GuildID)
	if err != nil {
		logging.Warn(err)
		return
	}

	guild.VerifiedRole = role.ID
	err = guild.Update()
	if err != nil {
		logging.Warn(err)
		return
	}

	utils.HiddenResponse(session, interaction, fmt.Sprintf("The verified role is now set to %v.", role.Mention()))
}

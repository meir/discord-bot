package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/discord-bot/internal/logging"
	"github.com/meir/discord-bot/pkg/structs"
	"go.mongodb.org/mongo-driver/mongo"
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
	if g, err := session.Guild(interaction.GuildID); err != nil || g.OwnerID != interaction.User.ID {
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   1 << 6,
				Content: "You have no permission for this command",
			},
		})
		return
	}
	role := interaction.ApplicationCommandData().Options[0].RoleValue(session, interaction.GuildID)

	guild, err := structs.NewQuery(session, interaction, db).Guild(interaction.GuildID)
	if err != nil {
		logging.Warn(err)
		return
	}

	guild.ModRole = role.ID
	err = guild.Update()
	if err != nil {
		logging.Warn(err)
		return
	}

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   1 << 6,
			Content: fmt.Sprintf("The moderator role is now set to %v.", role.Mention()),
		},
	})
}

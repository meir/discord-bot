package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/discord-bot/internal/logging"
	"github.com/meir/discord-bot/internal/utils"
	"github.com/meir/discord-bot/pkg/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	appendComponent(discordgo.ButtonComponent, "verification_request_button", verification_request_button)
}

/**

permissions

db.getchannel X

if db.channel = exists {
	if dsc.channel != exists {
		dsc.channel.create
		db.channel.update
		response
		return
	}
	response
	return
}
db.channel.create
dsc.channel.create
response
**/

func verification_request_button(session *discordgo.Session, interaction *discordgo.InteractionCreate, db *mongo.Database) {
	query := structs.NewQuery(session, interaction, db)

	guild, err := query.Guild(interaction.GuildID)
	if err != nil {
		logging.Warn(err)
		return
	}

	overwrites := []*discordgo.PermissionOverwrite{
		{
			ID:    interaction.Member.User.ID,
			Type:  discordgo.PermissionOverwriteTypeMember,
			Allow: discordgo.PermissionSendMessages | discordgo.PermissionViewChannel,
		},
		{
			ID:   interaction.GuildID, // > The @everyone role has the same ID as the guild it belongs to.
			Type: discordgo.PermissionOverwriteTypeRole,
			Deny: discordgo.PermissionAll,
		},
		{
			ID:    guild.ModRole,
			Type:  discordgo.PermissionOverwriteTypeRole,
			Allow: discordgo.PermissionSendMessages | discordgo.PermissionViewChannel,
		},
	}

	channel, err := query.ChannelByFilter(bson.M{
		"guild_id": interaction.GuildID,
		"metadata": bson.M{
			"type": structs.VERIFICATION_CHANNEL,
			"user": interaction.Member.User.ID,
		},
	})
	if err != nil && err != mongo.ErrNoDocuments {
		logging.Warn(err)
		return
	}

	if err != mongo.ErrNoDocuments {
		channels, err := session.GuildChannels(interaction.GuildID)
		if err != nil {
			logging.Warn(err)
			return
		}
		exists := false
		for _, v := range channels {
			if v.ID == channel.ChannelID {
				exists = true
			}
		}
		if !exists {
			ch, err := session.GuildChannelCreateComplex(interaction.GuildID, discordgo.GuildChannelCreateData{
				Name:                 fmt.Sprintf("verification-%v", interaction.Member.User.ID),
				Type:                 discordgo.ChannelTypeGuildText,
				Topic:                "You can explain here why you should be allowed into the server! Good luck :)",
				PermissionOverwrites: overwrites,
				ParentID:             guild.VerificationMessage.ParentID,
			})
			if err != nil {
				logging.Warn("failed to create channel", err)
				return
			}
			channel.ChannelID = ch.ID
			err = channel.Update()
			if err != nil {
				logging.Warn("failed to create channel", err)
				return
			}

			utils.HiddenResponse(session, interaction, fmt.Sprintf("Your verification channel has been created: %v", ch.Mention()))
			return
		}

		utils.HiddenResponse(session, interaction, fmt.Sprintf("You already have an open verification channel: <@%v>", channel.ChannelID))
		return
	}

	ch, err := session.GuildChannelCreateComplex(interaction.GuildID, discordgo.GuildChannelCreateData{
		Name:                 fmt.Sprintf("verification-%v", interaction.Member.User.ID),
		Type:                 discordgo.ChannelTypeGuildText,
		Topic:                "You can explain here why you should be allowed into the server! Good luck :)",
		PermissionOverwrites: overwrites,
		ParentID:             guild.VerificationMessage.ParentID,
	})
	if err != nil {
		logging.Warn("failed to create channel", err)
		return
	}

	channel = query.NewChannel(interaction.GuildID, ch.ID)
	channel.Metadata = map[string]string{
		"type": string(structs.VERIFICATION_CHANNEL),
		"user": interaction.Member.User.ID,
	}
	err = channel.Update()
	if err != nil {
		logging.Warn(err)
		return
	}

	utils.HiddenResponse(session, interaction, fmt.Sprintf("Your verification channel has been created: %v", ch.Mention()))
}

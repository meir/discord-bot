package events

import "github.com/bwmarrin/discordgo"

func Register(session *discordgo.Session) {
	session.AddHandler(OnReady)
}

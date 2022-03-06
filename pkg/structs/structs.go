package structs

type Guild struct {
	GuildID             string  `bson:"guild_id"`
	VerifiedRole        string  `bson:"verified_role"`
	VerificationMessage Message `bson:"verification_message"`
}

type Message struct {
	ChannelID string `bson:"channel_id"`
	MessageID string `bson:"message_id"`
	ParentID  string `bson:"parent_id"`
}

type Channel struct {
	GuildID   string            `bson:"guild_id"`
	ChannelID string            `bson:"channel_id"`
	Metadata  map[string]string `bson:"metadata"`
}

type ChannelType string

const (
	VERIFICATION_CHANNEL ChannelType = "verification_channel"
)

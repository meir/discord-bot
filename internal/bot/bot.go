package bot

import (
	"context"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/discord-bot/internal/commands"
	"github.com/meir/discord-bot/internal/logging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DiscordBot struct {
	Session  *discordgo.Session
	Database *mongo.Database
}

func NewDiscordBot() *DiscordBot {
	session, err := discordgo.New(fmt.Sprintf("Bot %v", os.Getenv("DISCORD_TOKEN")))

	if err != nil {
		logging.Fatal(err)
	}

	opts := options.Client()
	opts.ApplyURI(os.Getenv("MONGODB_URL"))
	opts.SetMaxPoolSize(5)

	var database *mongo.Client = nil
	if database, err = mongo.Connect(context.Background(), opts); err != nil {
		logging.Fatal(err)
	}

	err = database.Ping(context.Background(), nil)
	if err != nil {
		logging.Fatal(err)
	}

	return &DiscordBot{
		Session:  session,
		Database: database.Database(os.Getenv("MONGODB_DATABASE"), nil),
	}
}

func (d *DiscordBot) RegisterEvents() {
	d.Session.AddHandler(commands.Handler(d.Database))
}

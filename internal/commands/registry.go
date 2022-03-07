package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/meir/discord-bot/internal/logging"
	"go.mongodb.org/mongo-driver/mongo"
)

type Executor func(*discordgo.Session, *discordgo.InteractionCreate, *mongo.Database)

type Command struct {
	Name        string
	Command     string
	Description string
	Arguments   []*discordgo.ApplicationCommandOption
	Executor    Executor
}

type Component struct {
	Type     discordgo.ComponentType
	ID       string
	Executor Executor
}

var commands = map[string]Command{}
var components = map[string]Component{}
var registered = map[string]*discordgo.ApplicationCommand{}

func appendCommand(name, command, desc string, args []*discordgo.ApplicationCommandOption, e Executor) {
	commands[command] = Command{
		name, command, desc, args, e,
	}
}

func appendComponent(t discordgo.ComponentType, id string, e Executor) {
	components[id] = Component{
		t, id, e,
	}
}

func Handler(database *mongo.Database) interface{} {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commands[i.ApplicationCommandData().Name]; ok {
				h.Executor(s, i, database)
			}
		case discordgo.InteractionMessageComponent:
			if c, ok := components[i.MessageComponentData().CustomID]; ok {
				if i.MessageComponentData().ComponentType == c.Type {
					c.Executor(s, i, database)
				}
			}
		}
	}
}

func RegisterCommands(s *discordgo.Session) {
	logging.Println("Registering commands...")
	for _, g := range s.State.Guilds {
		for k, c := range commands {
			id, err := s.ApplicationCommandCreate(s.State.User.ID, g.ID, &discordgo.ApplicationCommand{
				Name:        k,
				Description: c.Description,
				Options:     c.Arguments,
			})
			if err != nil {
				logging.Warn("Cannot create '%v' command: %v", c.Name, err)
			}
			registered[c.Name] = id
		}
	}
}

func RemoveCommands(s *discordgo.Session) {
	for _, g := range s.State.Guilds {
		for _, c := range registered {
			s.ApplicationCommandDelete(s.State.User.ID, g.ID, c.ID)
		}
	}
}

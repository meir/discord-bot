package commands

// func init() {
// 	appendCommand("Reminder", "reminder", "Set a reminder to remind you once or with intervals", []*discordgo.ApplicationCommandOption{
// 		{
// 			Name:        "description",
// 			Description: "Describe what the reminder is for",
// 			Type:        discordgo.ApplicationCommandOptionString,
// 			Required:    true,
// 		},
// 		{
// 			Name:        "Repeating",
// 			Description: "Is the reminder for one time or repeating?",
// 			Type:        discordgo.ApplicationCommandOptionBoolean,
// 			Required:    true,
// 		},
// 		{
// 			Name:        "Interval",
// 			Description: "Interval for the reminder, ex: '5 minutes', '1 month', 'monthly', 'yearly'",
// 			Type:        discordgo.ApplicationCommandOptionString,
// 			Required:    true,
// 		},
// 	}, reminder)
// }

// func reminder(session *discordgo.Session, interaction *discordgo.InteractionCreate, database *mongo.Database) {

// 	type DurationFunc func(time.Time, int) time.Time

// 	keywords := map[string]DurationFunc{
// 		"second": func(t time.Time, mult int) time.Time {
// 			return t.Add(time.Second * time.Duration(mult))
// 		},
// 		"minute": func(t time.Time, mult int) time.Time {
// 			return t.Add(time.Minute * time.Duration(mult))
// 		},
// 		"hour": func(t time.Time, mult int) time.Time {
// 			return t.Add(time.Hour * time.Duration(mult))
// 		},
// 		"day": func(t time.Time, mult int) time.Time {
// 			return t.AddDate(0, 0, mult)
// 		},
// 		"week": func(t time.Time, mult int) time.Time {
// 			return t.AddDate(0, 0, 7*mult)
// 		},
// 		"month": func(t time.Time, mult int) time.Time {
// 			return t.AddDate(0, mult, 0)
// 		},
// 		"year": func(t time.Time, mult int) time.Time {
// 			return t.AddDate(mult, 0, 0)
// 		},
// 	}

// 	description := interaction.ApplicationCommandData().Options[0].StringValue()
// 	repeating := interaction.ApplicationCommandData().Options[1].BoolValue()
// 	interval := interaction.ApplicationCommandData().Options[2].StringValue()

// 	regex := regexp.MustCompile("^[0-9]{0,3} (second|minute|hour|day|week|month|year)[s]{0,1}$")
// 	if !regex.Match([]byte(interval)) {
// 		utils.HiddenResponse(session, interaction, "cant decipher the interval, needs to be something like '10 minutes', similar to 'x duration'.")
// 		return
// 	}

// 	reminder := structs.Reminder{}
// }

package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-errors/errors"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

var typeColor = map[string]int{
	"Info":    0x8af1ff,
	"Warning": 0xffff8a,
	"Fatal":   0xff948a,
	"Debug":   0xffffff,
}

func Println(l ...interface{}) {
	print(color.HiCyanString("[INFO]"), l...)
	webhook("Info", l...)
}

func Warn(l ...interface{}) {
	print(color.HiYellowString("[WARN]"), l...)
	webhook("Warning", l...)
}

func Fatal(l ...interface{}) {
	print(color.HiRedString("[FATAL]"), l...)
	webhook("Fatal", l...)
	panic(l)
}

func Debug(l ...interface{}) {
	print(color.HiWhiteString("[DEBUG]"), l...)
	webhook("Debug", l...)
}

func print(pre string, l ...interface{}) {
	l = append([]interface{}{pre}, l...)
	fmt.Println(l...)
}

func webhook(t string, l ...interface{}) {
	if os.Getenv("DEBUG_WEBHOOK") == "" {
		return
	}

	traces := []string{}
	for _, v := range l {
		if err, ok := v.(error); ok {
			trace := err.(*errors.Error).ErrorStack()
			traces = append(traces, trace)
			fmt.Println(trace)
		}
	}

	wh := discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       t,
				Color:       typeColor[t],
				Description: fmt.Sprintf("> %v\n%v", fmt.Sprintln(l...), strings.Join(traces, "\n---\n")),
				Timestamp:   time.Now().Format("2006-01-02 15:04"),
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("Build Version VH-%s", os.Getenv("VERSION")),
				},
			},
		},
	}

	data, err := json.Marshal(wh)
	if err != nil {
		fmt.Println("failed to marshal webhook params: ", err)
		return
	}

	r, err := http.Post(os.Getenv("DEBUG_WEBHOOK"), "application/json", bytes.NewReader(data))
	if err != nil {
		fmt.Println("failed to call debug webhook: ", err)
		return
	}
	defer r.Body.Close()
}

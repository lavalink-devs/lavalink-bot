package commands

import (
	"fmt"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/disgolink"
)

var info = discord.SlashCommandCreate{
	Name:        "info",
	Description: "Shows info about this bot",
}

func (c *Cmds) Info(e *handler.CommandEvent) error {
	return e.CreateMessage(discord.MessageCreate{
		Embeds: []discord.Embed{
			{
				Title:       "Lavalink Bot",
				Description: fmt.Sprintf("This bot is running on [disgo](https://github.com/disgoorg/disgo) `%s` and [disgolink](https://github.com/disgoorg/disgolink) `%s`.\n", disgo.Version, disgolink.Version),
			},
		},
	})
}

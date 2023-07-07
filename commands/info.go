package commands

import (
	"context"
	"fmt"
	"github.com/disgoorg/disgolink/v3/disgolink"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var info = discord.SlashCommandCreate{
	Name:        "info",
	Description: "Shows info about this bot",
}

func (c *Commands) Info(e *handler.CommandEvent) error {
	var fields []discord.EmbedField
	c.Lavalink.ForNodes(func(node disgolink.Node) {
		version, err := node.Version(context.TODO())
		var versionString string
		if err != nil {
			versionString = err.Error()
		} else {
			versionString = version
		}
		fields = append(fields, discord.EmbedField{
			Name:  node.Config().Name,
			Value: fmt.Sprintf("`%s`", versionString),
		})
	})
	return e.CreateMessage(discord.MessageCreate{
		Embeds: []discord.Embed{
			{
				Title: "Lavalink Bot",
				Fields: append([]discord.EmbedField{
					{
						Name:  "Source",
						Value: "[GitHub](https://github.com/lavalink-devs/lavalink-bot)",
					},
					{
						Name:  "DisGo",
						Value: fmt.Sprintf("`%s`", disgo.Version),
					},
					{
						Name:  "DisGoLink",
						Value: fmt.Sprintf("`%s`", disgolink.Version),
					},
				}, fields...),
			},
		},
	})
}

package commands

import (
	"context"
	"fmt"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/json"
)

var info = discord.SlashCommandCreate{
	Name:        "info",
	Description: "Shows info about this bot",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "bot",
			Description: "Shows info about this bot",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "lavalink",
			Description: "Shows info about the lavalink nodes",
		},
	},
}

func (c *Commands) InfoBot(e *handler.CommandEvent) error {
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

func (c *Commands) InfoLavalink(e *handler.CommandEvent) error {
	nodeInfos := map[string]lavalink.Info{}
	c.Lavalink.ForNodes(func(node disgolink.Node) {
		nodeInfo, err := node.Info(context.TODO())
		if err != nil {
			return
		}
		nodeInfos[node.Config().Name] = *nodeInfo
	})

	rawInfo, err := json.MarshalIndent(nodeInfos, "", "  ")
	if err != nil {
		return e.CreateMessage(discord.MessageCreate{
			Content: "Failed to marshal lavalink info: " + err.Error(),
		})
	}

	return e.CreateMessage(discord.MessageCreate{
		Embeds: []discord.Embed{
			{
				Title:       "Lavalink Nodes",
				Description: "```json\n" + string(rawInfo) + "\n```",
			},
		},
	})
}

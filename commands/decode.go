package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/json"
)

var decode = discord.SlashCommandCreate{
	Name:        "decode",
	Description: "Decode a base64 encoded lavalink track",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "track",
			Description: "The base64 encoded lavalink track",
			Required:    true,
		},
	},
}

func (c *Cmds) Decode(e *handler.CommandEvent) error {
	track := e.SlashCommandInteractionData().String("track")

	if err := e.DeferCreateMessage(false); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	decoded, err := c.Lavalink.BestNode().Rest().DecodeTracks(ctx, []string{track})
	if err != nil {
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: json.Ptr(fmt.Sprintf("failed to decode track: %s", err)),
		})
		return err
	}

	data, err := json.MarshalIndent(decoded[0], "", "  ")
	if err != nil {
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: json.Ptr(fmt.Sprintf("failed to decode track: %s", err)),
		})
		return err
	}

	_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
		Content: json.Ptr(fmt.Sprintf("```json\n%s\n```", data)),
	})
	return err
}

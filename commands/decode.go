package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/json"

	"github.com/lavalink-devs/lavalink-bot/internal/trackdecode"
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
		discord.ApplicationCommandOptionBool{
			Name:        "lavalink",
			Description: "If the Lavalink server should be used for decoding",
			Required:    false,
		},
	},
}

func (c *Commands) Decode(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	track := data.String("track")
	lavalink := data.Bool("lavalink")

	if !lavalink {
		var content string
		decoded, version, err := trackdecode.DecodeString(track)
		if err != nil {
			content += fmt.Sprintf("error while decoding track: %s\n", err.Error())
		}
		if version > 0 {
			content += fmt.Sprintf("track was encoded with version: `%d`\n", version)
		}
		if decoded != nil {
			decodedData, _ := json.MarshalIndent(decoded, "", "  ")
			content += fmt.Sprintf("```json\n%s\n```", decodedData)
		}

		return e.CreateMessage(discord.MessageCreate{
			Content: content,
		})
	}

	if err := e.DeferCreateMessage(false); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(e.Ctx, 10*time.Second)
	defer cancel()
	decoded, err := c.Lavalink.BestNode().Rest().DecodeTrack(ctx, track)
	if err != nil {
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: json.Ptr(fmt.Sprintf("failed to decode track: %s", err)),
		})
		return err
	}

	decodedData, err := json.MarshalIndent(decoded, "", "  ")
	if err != nil {
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: json.Ptr(fmt.Sprintf("failed to decode track: %s", err)),
		})
		return err
	}

	_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
		Content: json.Ptr(fmt.Sprintf("```json\n%s\n```", decodedData)),
	})
	return err
}

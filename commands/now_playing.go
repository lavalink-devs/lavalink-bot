package commands

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/json"
	"github.com/lavalink-devs/lavalink-bot/internal/res"
)

func (c *Commands) NowPlaying(e *handler.CommandEvent) error {
	player := c.Lavalink.Player(*e.GuildID())
	track := player.Track()
	if track == nil {
		return e.CreateMessage(discord.MessageCreate{
			Content: "There is no track playing.",
			Flags:   discord.MessageFlagEphemeral,
		})
	}
	content := fmt.Sprintf("Now playing: %s", res.FormatTrack(*track, player.Position()))

	if e.SlashCommandInteractionData().Bool("raw") {
		data, err := json.MarshalIndent(track, "", "  ")
		if err != nil {
			return e.CreateMessage(discord.MessageCreate{
				Content: fmt.Sprintf("Failed to marshal track: %s", err),
				Flags:   discord.MessageFlagEphemeral,
			})
		}

		content += fmt.Sprintf("\n```json\n%s\n```", data)
	}

	return e.CreateMessage(discord.MessageCreate{
		Content: content,
	})
}

package commands

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"

	"github.com/lavalink-devs/lavalink-bot/internal/res"
)

func (c *Commands) Queue(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	_, tracks := c.MusicQueue.Get(*e.GuildID())
	if len(tracks) == 0 {
		return e.CreateMessage(discord.MessageCreate{
			Content: "No tracks in queue",
		})
	}

	content := fmt.Sprintf("**Queue(%d):**\n", len(tracks))
	for i, track := range tracks {
		line := fmt.Sprintf("%d. %s\n", i+1, res.FormatTrack(track, 0))
		if len(content)+len(line) > 1980 {
			content += fmt.Sprintf("... and %d more", len(tracks)-i)
			break
		}
		content += line
	}

	return e.CreateMessage(discord.MessageCreate{
		Content: content,
	})
}

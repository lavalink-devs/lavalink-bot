package commands

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/lavalink-devs/lavalink-bot/internal/res"
)

func (c *Commands) Queue(e *handler.CommandEvent) error {
	_, tracks := c.MusicQueue.Get(*e.GuildID())
	if len(tracks) == 0 {
		return e.CreateMessage(discord.MessageCreate{
			Content: fmt.Sprintf("No tracks in queue"),
		})
	}

	content := fmt.Sprintf("**Queue(%d):**\n", len(tracks))
	for i, track := range tracks {
		content += fmt.Sprintf("%d. %s\n", i+1, res.FormatTrack(track, 0))
	}

	return e.CreateMessage(discord.MessageCreate{
		Content: content,
	})
}

package commands

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func (c *Commands) Remove(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	index := data.Int("index")
	ok := c.MusicQueue.Remove(*e.GuildID(), index-1, index)
	if !ok {
		return e.CreateMessage(discord.MessageCreate{
			Content: fmt.Sprintf("Failed to remove track %d from queue", index),
		})
	}

	return e.CreateMessage(discord.MessageCreate{
		Content: fmt.Sprintf("Removed track %d from queue", index),
	})
}

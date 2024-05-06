package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func (c *Commands) Shuffle(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	ok := c.MusicQueue.Shuffle(*e.GuildID())
	if !ok {
		return e.CreateMessage(discord.MessageCreate{
			Content: "No or not enough tracks in queue to shuffle",
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	return e.CreateMessage(discord.MessageCreate{
		Content: "🔀 Shuffled queue",
	})
}

package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func (c *Commands) Clear(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	_, tracks := c.MusicQueue.Get(*e.GuildID())
	if len(tracks) == 0 {
		return e.CreateMessage(discord.MessageCreate{
			Content: "Queue is already empty",
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	c.MusicQueue.Clear(*e.GuildID())

	return e.CreateMessage(discord.MessageCreate{
		Content: "🗑️ Cleared the queue",
	})
}

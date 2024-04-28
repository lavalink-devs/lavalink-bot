package commands

import (
	"context"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func (c *Commands) Disconnect(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	ctx, cancel := context.WithTimeout(e.Ctx, 10*time.Second)
	defer cancel()

	if err := c.Client.UpdateVoiceState(ctx, *e.GuildID(), nil, false, false); err != nil {
		return e.CreateMessage(discord.MessageCreate{
			Content: "Failed to disconnect player",
		})
	}
	return e.CreateMessage(discord.MessageCreate{
		Content: "Disconnected player",
	})
}

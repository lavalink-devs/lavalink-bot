package commands

import (
	"context"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func (c *Cmds) Disconnect(e *handler.CommandEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

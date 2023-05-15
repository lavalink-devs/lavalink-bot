package commands

import (
	"context"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

func (c *Cmds) Stop(e *handler.CommandEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	player := c.Lavalink.ExistingPlayer(*e.GuildID())
	if err := player.Update(ctx, lavalink.WithNullTrack()); err != nil {
		return e.CreateMessage(discord.MessageCreate{
			Content: "Failed to stop player",
		})
	}
	return e.CreateMessage(discord.MessageCreate{
		Content: "Stopped player",
	})
}

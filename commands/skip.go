package commands

import (
	"context"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

func (c *Commands) Skip(e *handler.CommandEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	player := c.Lavalink.ExistingPlayer(*e.GuildID())
	track, ok := c.MusicQueue.Next(*e.GuildID())
	if !ok {
		return e.CreateMessage(discord.MessageCreate{
			Content: "No more tracks in queue",
			Flags:   discord.MessageFlagEphemeral,
		})
	}
	if err := player.Update(ctx, lavalink.WithTrack(track)); err != nil {
		return e.CreateMessage(discord.MessageCreate{
			Content: "Failed to play skip track",
		})
	}

	return e.CreateMessage(discord.MessageCreate{
		Content: "⏭ Skipped track",
	})
}

package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

func (c *Commands) Skip(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	ctx, cancel := context.WithTimeout(e.Ctx, 10*time.Second)
	defer cancel()

	var trackCount = data.Int("count")
	if trackCount == 0 {
		trackCount = 1
	}

	player := c.Lavalink.ExistingPlayer(*e.GuildID())
	track, ok := c.MusicQueue.NextCount(*e.GuildID(), trackCount)
	if !ok {
		return e.CreateMessage(discord.MessageCreate{
			Content: "Not enough tracks to skip",
			Flags:   discord.MessageFlagEphemeral,
		})
	}
	if err := player.Update(ctx, lavalink.WithTrack(track)); err != nil {
		return e.CreateMessage(discord.MessageCreate{
			Content: "Failed to skip track(s)",
		})
	}

	return e.CreateMessage(discord.MessageCreate{
		Content: fmt.Sprintf("⏭ Skipped %d track(s)", trackCount),
	})
}

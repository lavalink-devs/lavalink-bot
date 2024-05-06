package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/lavalink-devs/lavalink-bot/internal/res"
)

func (c *Commands) Seek(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	ctx, cancel := context.WithTimeout(e.Ctx, 10*time.Second)
	defer cancel()
	player := c.Lavalink.ExistingPlayer(*e.GuildID())

	position := data.Int("position")
	duration, ok := data.OptInt("unit")
	if !ok {
		duration = int(lavalink.Second)
	}

	newPosition := lavalink.Duration(position * duration)
	if err := player.Update(ctx, lavalink.WithPosition(newPosition)); err != nil {
		return e.CreateMessage(discord.MessageCreate{
			Content: "Failed to seek to position",
		})
	}

	return e.CreateMessage(discord.MessageCreate{
		Content: fmt.Sprintf("⏩ Seeked to %s", res.FormatDuration(newPosition)),
	})
}

package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

func (c *Commands) Volume(e *handler.CommandEvent) error {
	volume := e.SlashCommandInteractionData().Int("volume")
	player := c.Lavalink.ExistingPlayer(*e.GuildID())
	oldVolume := player.Volume()

	if volume == oldVolume {
		return e.CreateMessage(discord.MessageCreate{
			Content: fmt.Sprintf("Volume is already at `%d%%`", volume),
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := player.Update(ctx, lavalink.WithVolume(volume)); err != nil {
		return e.CreateMessage(discord.MessageCreate{
			Content: "Failed to set volume",
		})
	}

	var content string
	if volume > oldVolume {
		content = fmt.Sprintf("🔊 Increased volume to `%d%%`", volume)
	} else {
		content = fmt.Sprintf("🔉 Decreased volume to `%d%%`", volume)
	}

	return e.CreateMessage(discord.MessageCreate{
		Content: content,
	})
}

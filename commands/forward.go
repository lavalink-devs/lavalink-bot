package commands

import (
	"context"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/lavalink"

	"github.com/lavalink-devs/lavalink-bot/internal/res"
)

func (c *Commands) Forward(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	ctx, cancel := context.WithTimeout(e.Ctx, 10*time.Second)
	defer cancel()
	player := c.Lavalink.ExistingPlayer(*e.GuildID())

	amt := data.Int("amount")
	unit, ok := data.OptInt("unit")
	if !ok {
		unit = int(lavalink.Second)
	}

	cPos := int(player.State().Position)
	nPos := cPos + (amt * unit)

	if nPos >= int(player.Track().Info.Length) {
		nPos = int(player.Track().Info.Length) - 1
	}

	if nPos == cPos {
		return e.CreateMessage(discord.MessageCreate{
			Content: "The player is already at position **`" + res.FormatDuration(lavalink.Duration(cPos)) + "`**",
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	if err := player.Update(ctx, lavalink.WithPosition(lavalink.Duration(nPos))); err != nil {
		return e.CreateMessage(discord.MessageCreate{
			Content: "Failed to forward the track.",
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	return e.CreateMessage(discord.MessageCreate{
		Content: "⏩ Forwarded to **`" + res.FormatDuration(lavalink.Duration(nPos)) + "`**",
	})
}

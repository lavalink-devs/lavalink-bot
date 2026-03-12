package commands

import (
	"context"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v4/disgolink"
)

func (c *Commands) Stop(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	ctx, cancel := context.WithTimeout(e.Ctx, 10*time.Second)
	defer cancel()
	player := c.Lavalink.ExistingPlayer(*e.GuildID())
	if err := player.Update(ctx, disgolink.WithNullTrack()); err != nil {
		return e.CreateMessage(discord.MessageCreate{
			Content: "Failed to stop player",
		})
	}
	return e.CreateMessage(discord.MessageCreate{
		Content: "Stopped player",
	})
}

package commands

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/lavalink-devs/lavalink-bot/internal/res"
)

func (c *Commands) NowPlaying(e *handler.CommandEvent) error {
	player := c.Lavalink.Player(*e.GuildID())
	track := player.Track()
	if track == nil {
		return e.CreateMessage(discord.MessageCreate{
			Content: "There is no track playing.",
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	return e.CreateMessage(discord.MessageCreate{
		Content: fmt.Sprintf("Now playing: %s", res.FormatTrack(*track, player.Position())),
	})
}

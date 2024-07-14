package commands

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"

	"github.com/lavalink-devs/lavalink-bot/lavalinkbot"
)

func (c *Commands) Repeat(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	newMode := lavalinkbot.RepeatMode(data.String("mode"))

	c.MusicQueue.SetRepeatMode(*e.GuildID(), newMode)

	return e.CreateMessage(discord.MessageCreate{
		Content: fmt.Sprintf("Set repeat mode to `%s`", newMode),
	})
}

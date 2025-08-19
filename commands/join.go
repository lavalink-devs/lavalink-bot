package commands

import (
	"context"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
)

func (c *Commands) Join(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	var vcID snowflake.ID

	ch, ok := data.OptChannel("channel")
	if ok && ch.Type == discord.ChannelTypeGuildVoice {
		vcID = ch.ID
	} else {
		vs, ok := c.Client.Caches().VoiceState(*e.GuildID(), e.User().ID)
		if !ok || vs.ChannelID == nil {
			return e.CreateMessage(discord.MessageCreate{
				Content: "You must be in a voice channel or mention one",
				Flags:   discord.MessageFlagEphemeral,
			})
		}
		vcID = *vs.ChannelID
	}

	if err := c.Client.UpdateVoiceState(context.Background(), *e.GuildID(), &vcID, false, false); err != nil {
		return e.CreateMessage(discord.MessageCreate{
			Content: fmt.Sprintf("Failed to join voice channel: %s", err),
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	return e.CreateMessage(discord.MessageCreate{
		Content: fmt.Sprintf("Joined <#%s>", vcID.String()),
	})
}

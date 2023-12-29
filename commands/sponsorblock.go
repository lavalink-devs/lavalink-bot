package commands

import (
	"context"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/sponsorblock-plugin"
)

func (c *Commands) ShowSponsorblock(e *handler.CommandEvent) error {
	node := c.Lavalink.Player(*e.GuildID()).Node()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	categories, err := sponsorblock.GetCategories(ctx, node.Rest(), node.SessionID(), *e.GuildID())
	if err != nil {
		return e.CreateMessage(discord.MessageCreate{
			Content: "Failed to get categories: " + err.Error(),
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	content := "Sponsorblock categories:\n"
	for _, category := range categories {
		content += "* " + string(category) + "\n"
	}

	return e.CreateMessage(discord.MessageCreate{
		Content: content,
	})
}

func (c *Commands) SetSponsorblock(e *handler.CommandEvent) error {
	data := e.SlashCommandInteractionData()
	node := c.Lavalink.Player(*e.GuildID()).Node()

	var categories []sponsorblock.SegmentCategory
	for _, category := range sponsorblockOptions {
		if data.Bool(category.OptionName()) {
			categories = append(categories, sponsorblock.SegmentCategory(category.OptionName()))
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := sponsorblock.SetCategories(ctx, node.Rest(), node.SessionID(), *e.GuildID(), categories); err != nil {
		return e.CreateMessage(discord.MessageCreate{
			Content: "Failed to set categories: " + err.Error(),
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	return e.CreateMessage(discord.MessageCreate{
		Content: "Sponsorblock categories set",
	})
}

package commands

import (
	"context"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/lavalink-devs/lavalink-bot/lavalinkbot"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

var latest = discord.SlashCommandCreate{
	Name:        "latest",
	Description: "Shows the latest version of Lavalink and Plugins",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:         "type",
			Description:  "The type of version you want to see",
			Required:     false,
			Autocomplete: true,
		},
	},
}

func (c *Commands) LatestAutocomplete(e *handler.AutocompleteEvent) error {
	options := []string{"lavalink"}
	for _, plugin := range c.Cfg.Plugins {
		options = append(options, plugin.Dependency)
	}

	ranks := fuzzy.RankFindFold(e.Data.String("type"), options)

	var choices []discord.AutocompleteChoice
	for _, rank := range ranks {
		if len(choices) >= 25 {
			break
		}
		choices = append(choices, discord.AutocompleteChoiceString{
			Name:  options[rank.OriginalIndex],
			Value: options[rank.OriginalIndex],
		})
	}

	return e.AutocompleteResult(choices)
}

func (c *Commands) Latest(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	if err := e.DeferCreateMessage(false); err != nil {
		return err
	}

	latestType, ok := e.SlashCommandInteractionData().OptString("type")
	if !ok {
		embed := getLatestEmbed(c, "lavalink")
		for _, plugin := range c.Cfg.Plugins {
			newEmbed := getLatestEmbed(c, plugin.Dependency)
			embed.Fields = append(embed.Fields, discord.EmbedField{
				Name:  newEmbed.Title,
				Value: newEmbed.Description,
			})
		}

		_, err := e.UpdateInteractionResponse(discord.MessageUpdate{
			Embeds: &[]discord.Embed{embed},
		})
		return err
	}
	embed := getLatestEmbed(c, latestType)
	_, err := e.UpdateInteractionResponse(discord.MessageUpdate{
		Embeds: &[]discord.Embed{embed},
	})
	return err
}

func getLatestEmbed(c *Commands, latestType string) discord.Embed {
	if latestType == "lavalink" {
		embed := discord.Embed{
			Author: &discord.EmbedAuthor{
				Name: "Latest Lavalink Version",
			},
		}
		release, _, err := c.GitHub.Repositories.GetLatestRelease(context.Background(), "lavalink-devs", "Lavalink")
		if err != nil {
			embed.Description = "Failed to get latest Lavalink version: " + err.Error()
			return embed
		}
		embed.Author.URL = release.GetHTMLURL()
		embed.Description = fmt.Sprintf("**Version:** `%s`\n**Release Date:** %s", release.GetTagName(), discord.NewTimestamp(discord.TimestampStyleLongDateTime, release.GetPublishedAt().Time))
		return embed
	}
	var pluginCfg lavalinkbot.PluginConfig
	for _, plugin := range c.Cfg.Plugins {
		if plugin.Dependency == latestType {
			pluginCfg = plugin
			break
		}
	}
	if pluginCfg.Dependency == "" {
		return discord.Embed{
			Description: fmt.Sprintf("Unknown plugin: `%s`", latestType),
		}
	}

	metadata, err := c.Maven.FetchLatestVersion(pluginCfg.Dependency, pluginCfg.Repository)
	if err != nil {
		return discord.Embed{
			Description: fmt.Sprintf("Failed to get latest %s version: %s", pluginCfg.Dependency, err.Error()),
		}
	}
	return discord.Embed{
		Title:       fmt.Sprintf("Latest %s Version", metadata.ArtifactID),
		Description: fmt.Sprintf("**Version:** `%s`", metadata.Latest()),
	}

}

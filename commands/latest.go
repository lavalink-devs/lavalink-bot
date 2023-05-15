package commands

import (
	"context"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/lavalink-devs/lavalink-bot/internal/maven"
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

func (c *Cmds) LatestAutocomplete(e *handler.AutocompleteEvent) error {
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

	return e.Result(choices)
}

func (c *Cmds) Latest(e *handler.CommandEvent) error {
	if err := e.DeferCreateMessage(false); err != nil {
		return err
	}

	latestType, ok := e.SlashCommandInteractionData().OptString("type")
	if !ok {
		embed, err := getLatestEmbed(c, "lavalink")
		if err != nil {
			return err
		}
		for _, plugin := range c.Cfg.Plugins {
			newEmbed, err := getLatestEmbed(c, plugin.Dependency)
			if err != nil {
				return err
			}
			embed.Fields = append(embed.Fields, discord.EmbedField{
				Name:  newEmbed.Title,
				Value: newEmbed.Description,
			})
		}

		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Embeds: &[]discord.Embed{*embed},
		})
		return err
	}
	embed, err := getLatestEmbed(c, latestType)
	if err != nil {
		return err
	}
	_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
		Embeds: &[]discord.Embed{*embed},
	})
	return err
}

func getLatestEmbed(c *Cmds, latestType string) (*discord.Embed, error) {
	if latestType == "lavalink" {
		release, _, err := c.GitHub.Repositories.GetLatestRelease(context.Background(), "lavalink-devs", "Lavalink")
		if err != nil {
			return nil, err
		}
		return &discord.Embed{
			Author: &discord.EmbedAuthor{
				Name: "Latest Lavalink Version",
				URL:  release.GetHTMLURL(),
			},
			Description: fmt.Sprintf("**Version:** `%s`\n**Release Date:** %s", release.GetTagName(), discord.NewTimestamp(discord.TimestampStyleLongDateTime, release.GetPublishedAt().Time)),
		}, nil
	}
	var pluginCfg lavalinkbot.PluginConfig
	for _, plugin := range c.Cfg.Plugins {
		if plugin.Dependency == latestType {
			pluginCfg = plugin
			break
		}
	}
	if pluginCfg.Dependency == "" {
		return nil, fmt.Errorf("plugin %s not found", latestType)
	}

	metadata, err := maven.FetchLatestVersion(c.HTTPClient, pluginCfg.Dependency, pluginCfg.Repository)
	if err != nil {
		return nil, err
	}
	return &discord.Embed{
		Title:       fmt.Sprintf("Latest %s Version", metadata.ArtifactID),
		Description: fmt.Sprintf("**Version:** `%s`", metadata.Latest()),
	}, nil

}

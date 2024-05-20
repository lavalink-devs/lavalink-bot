package commands

import (
	"context"
	"fmt"
	"strings"
	"sync"

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
	options := []string{"Lavalink"}
	for _, plugin := range c.Cfg.Plugins {
		options = append(options, plugin.Name)
	}

	ranks := fuzzy.RankFindFold(e.Data.String("type"), options)

	var choices []discord.AutocompleteChoice
	for _, rank := range ranks {
		if len(choices) >= 25 {
			break
		}

		var value string
		if rank.OriginalIndex == 0 {
			value = "lavalink"
		} else {
			value = c.Cfg.Plugins[rank.OriginalIndex-1].Dependency
		}

		choices = append(choices, discord.AutocompleteChoiceString{
			Name:  options[rank.OriginalIndex],
			Value: value,
		})
	}

	return e.AutocompleteResult(choices)
}

func (c *Commands) Latest(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	if err := e.DeferCreateMessage(false); err != nil {
		return err
	}

	var types []string
	latestType, ok := data.OptString("type")
	if ok {
		types = []string{latestType}
	} else {
		types = append(types, "lavalink")
		for _, plugin := range c.Cfg.Plugins {
			types = append(types, plugin.Dependency)
		}
	}

	var wg sync.WaitGroup
	versions := make([]string, len(types))
	for i, versionType := range types {
		wg.Add(1)
		go func() {
			defer wg.Done()
			versions[i] = getLatest(c, versionType)
		}()
	}
	wg.Wait()

	_, err := e.UpdateInteractionResponse(discord.MessageUpdate{
		Embeds: &[]discord.Embed{
			{
				Title:       "Latest Versions",
				URL:         "https://lavalink.dev/plugins.html",
				Description: strings.Join(versions, "\n\n"),
			},
		},
	})
	return err
}

func getLatest(c *Commands, latestType string) string {
	if latestType == "lavalink" {
		release, _, err := c.GitHub.Repositories.GetLatestRelease(context.Background(), "lavalink-devs", "Lavalink")
		if err != nil {
			return "Failed to get latest Lavalink version: " + err.Error()
		}
		return fmt.Sprintf("**[`Lavalink`](%s):**\n**Latest Version:** `%s`\n**Release Date:** %s", release.GetHTMLURL(), release.GetTagName(), discord.NewTimestamp(discord.TimestampStyleLongDateTime, release.GetPublishedAt().Time))
	}

	var pluginCfg lavalinkbot.PluginConfig
	for _, plugin := range c.Cfg.Plugins {
		if plugin.Dependency == latestType {
			pluginCfg = plugin
			break
		}
	}
	if pluginCfg.Dependency == "" {
		return fmt.Sprintf("Unknown plugin: `%s`", latestType)
	}

	metadata, err := c.Maven.FetchLatestVersion(pluginCfg.Dependency, pluginCfg.Repository)
	if err != nil {
		return fmt.Sprintf("Failed to get latest %s version: %s", pluginCfg.Dependency, err.Error())
	}
	return fmt.Sprintf("**[`%s`](%s):**\n**Latest Version:** `%s`", pluginCfg.Name, pluginCfg.Git, metadata.Latest())
}

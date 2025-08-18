package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func (c *Commands) RemoveAutocomplete(e *handler.AutocompleteEvent) error {
	_, tracks := c.MusicQueue.Get(*e.GuildID())
	if len(tracks) == 0 {
		return e.AutocompleteResult(nil)
	}

	query := e.Data.String("index")
	choices := make([]discord.AutocompleteChoice, 0, len(tracks))
	for i, track := range tracks {
		name := fmt.Sprintf("%d: %s - %s", i+1, track.Info.Title, track.Info.Author)
		if len(name) > 100 {
			name = name[:97] + "..."
		}
		choices = append(choices, discord.AutocompleteChoiceString{
			Name:  name,
			Value: strconv.Itoa(i + 1),
		})
	}
	if query != "" {
		f := make([]discord.AutocompleteChoice, 0)
		q := strings.ToLower(query)
		for _, choice := range choices {
			choiceName := choice.(discord.AutocompleteChoiceString).Name
			if strings.Contains(strings.ToLower(choiceName), q) {
				f = append(f, choice)
			}
		}
		choices = f
	}
	if len(choices) > 25 {
		choices = choices[:25]
	}

	return e.AutocompleteResult(choices)
}

func (c *Commands) Remove(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	index := int(data.Int("index"))

	_, tracks := c.MusicQueue.Get(*e.GuildID())
	if len(tracks) == 0 {
		return e.CreateMessage(discord.MessageCreate{
			Content: "Queue is empty",
		})
	}

	if index < 1 || index > len(tracks) {
		return e.CreateMessage(discord.MessageCreate{
			Content: fmt.Sprintf("Index must be between 1 and %d", len(tracks)),
		})
	}

	c.MusicQueue.Remove(*e.GuildID(), index-1, index)

	return e.CreateMessage(discord.MessageCreate{
		Content: fmt.Sprintf("🗑️ Removed track from **#%d** queue", index),
	})
}

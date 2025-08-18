package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func (c *Commands) MoveAutocomplete(e *handler.AutocompleteEvent) error {
	_, tracks := c.MusicQueue.Get(*e.GuildID())
	if len(tracks) == 0 {
		return e.AutocompleteResult(nil)
	}

	query := e.Data.String("from")
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
			if strings.Contains(strings.ToLower(choice.(discord.AutocompleteChoiceString).Name), q) {
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

func (c *Commands) Move(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	from := data.Int("from")
	to := data.Int("to")

	_, tracks := c.MusicQueue.Get(*e.GuildID())
	if len(tracks) < 2 {
		return e.CreateMessage(discord.MessageCreate{
			Content: "Queue is empty",
		})
	}

	if from < 1 || from > len(tracks) || to < 1 || to > len(tracks) {
		return e.CreateMessage(discord.MessageCreate{
			Content: fmt.Sprintf("Index ust be between 1 and %d", len(tracks)),
		})
	}

	if from == to {
		return e.CreateMessage(discord.MessageCreate{
			Content: "Both indexes are same",
		})
	}

	c.MusicQueue.Move(*e.GuildID(), from-1, to-1)

	return e.CreateMessage(discord.MessageCreate{
		Content: fmt.Sprintf("Moved **%s** from position %d → %d", tracks[from-1].Info.Title, from, to),
	})
}

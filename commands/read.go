package commands

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
	"go.deanishe.net/fuzzy"

	"github.com/lavalink-devs/lavalink-bot/internal/res"
)

var read = discord.SlashCommandCreate{
	Name:        "read",
	Description: "Tells someone to read something",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:         "thing",
			Description:  "The thing someone should read",
			Required:     true,
			Autocomplete: true,
		},
		discord.ApplicationCommandOptionUser{
			Name:        "user",
			Description: "The user to tell to read something",
			Required:    false,
		},
	},
	Contexts: []discord.InteractionContextType{
		discord.InteractionContextTypeGuild,
		discord.InteractionContextTypeBotDM,
	},
}

func (c *Commands) Read(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	var msg discord.MessageCreate
	user, ok := data.OptUser("user")
	if ok {
		msg.Content += fmt.Sprintf("Hey %s,\n", user.Mention())
		msg.AllowedMentions = &discord.AllowedMentions{
			Users: []snowflake.ID{user.ID},
		}
	}

	thing, ok := c.Things[data.String("thing")]
	if !ok {
		return e.CreateMessage(discord.MessageCreate{
			Content: "I don't know that thing",
			Flags:   discord.MessageFlagEphemeral,
		})
	}
	msg.Content += thing.Content

	for _, file := range thing.Files {
		msg.Files = append(msg.Files, discord.NewFile(file.Name, "", file.Reader()))
	}

	return e.CreateMessage(msg)
}

func (c *Commands) ReadAutocomplete(e *handler.AutocompleteEvent) error {
	thing := e.Data.String("thing")

	var choices []discord.AutocompleteChoice
	for _, t := range c.Things {
		choices = append(choices, discord.AutocompleteChoiceString{
			Name:  res.Trim(t.Name, 100),
			Value: t.FileName,
		})
	}

	fuzzy.Sort(Choices(choices), thing)

	if len(choices) > 25 {
		choices = choices[:25]
	}

	return e.AutocompleteResult(choices)
}

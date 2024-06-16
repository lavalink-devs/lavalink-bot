package commands

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/json"
)

var resolve = discord.SlashCommandCreate{
	Name:        "resolve",
	Description: "Resolve an identifier to it's result",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "identifier",
			Description: "The identifier to resolve",
			Required:    true,
		},
	},
}

func (c *Commands) Resolve(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	identifier := data.String("identifier")

	if err := e.DeferCreateMessage(false); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(e.Ctx, 10*time.Second)
	defer cancel()
	result, err := c.Lavalink.BestNode().Rest().LoadTracks(ctx, identifier)
	if err != nil {
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: json.Ptr(fmt.Sprintf("failed to resolve identifier: %s", err)),
		})
		return err
	}

	var (
		content string
		files   []*discord.File
	)
	switch result.LoadType {
	case lavalink.LoadTypeTrack, lavalink.LoadTypePlaylist, lavalink.LoadTypeSearch:
		decodedData, err := json.MarshalIndent(result.Data, "", "  ")
		if err != nil {
			_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
				Content: json.Ptr(fmt.Sprintf("failed to resolve identifier: %s", err)),
			})
			return err
		}

		if len(decodedData) > 2000 {
			files = append(files, &discord.File{
				Name:   "result.json",
				Reader: bytes.NewReader(decodedData),
			})
			content = "result is too long, see attached file"
		} else {
			content = fmt.Sprintf("```json\n%s\n```", decodedData)
		}
		content = fmt.Sprintf("LoadType: `%s`\n%s", result.LoadType, content)

	case lavalink.LoadTypeEmpty:
		content = "LoadType: `empty`"
	case lavalink.LoadTypeError:
		ex, _ := result.Data.(lavalink.Exception)
		if ex.Cause != nil {
			files = append(files, &discord.File{
				Name:   "cause.txt",
				Reader: bytes.NewReader([]byte(*ex.Cause)),
			})
		}
		content = fmt.Sprintf("LoadType: `error`\nMessage: %s\nSeverity: %s", ex.Message, ex.Severity)
	}

	_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
		Content: &content,
		Files:   files,
	})
	return err
}

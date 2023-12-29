package commands

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/json"
	"github.com/disgoorg/lavalyrics-plugin"
)

func (c *Commands) Lyrics(e *handler.CommandEvent) error {
	if err := e.DeferCreateMessage(false); err != nil {
		return err
	}

	node := c.Lavalink.BestNode()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	lyrics, err := lavalyrics.GetLyrics(ctx, node.Rest(), node.SessionID(), *e.GuildID())
	if err != nil {
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: json.Ptr(fmt.Sprintf("failed to decode track: %s", err)),
		})
		return err
	}

	var content string
	if len(lyrics.Lines) == 0 {
		content = lyrics.Text
	} else {
		for _, line := range lyrics.Lines {
			content += fmt.Sprintf("%s\n", line.Line)
		}
	}

	_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
		Files: []*discord.File{
			discord.NewFile("lyrics.txt", "", bytes.NewReader([]byte(content))),
		},
	})
	return err
}

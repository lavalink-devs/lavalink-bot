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

	"github.com/lavalink-devs/lavalink-bot/internal/res"
)

func (c *Commands) Lyrics(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	skipTrackSource := data.Bool("skip-track-source")

	var (
		track  string
		lyrics *lavalyrics.Lyrics
		err    error
	)

	ctx, cancel := context.WithTimeout(e.Ctx, 10*time.Second)
	defer cancel()
	if encodedTrack, ok := data.OptString("track"); ok {
		track = fmt.Sprintf("`%s`", encodedTrack)

		if err = e.DeferCreateMessage(false); err != nil {
			return err
		}

		lyrics, err = lavalyrics.GetLyrics(ctx, c.Lavalink.BestNode().Rest(), encodedTrack, skipTrackSource)
	} else {
		player := c.Lavalink.ExistingPlayer(*e.GuildID())
		if player == nil {
			return e.CreateMessage(discord.MessageCreate{
				Content: "No player found",
				Flags:   discord.MessageFlagEphemeral,
			})
		}

		playingTrack := player.Track()
		if playingTrack == nil {
			return e.CreateMessage(discord.MessageCreate{
				Content: "no track playing",
				Flags:   discord.MessageFlagEphemeral,
			})
		}

		track = res.FormatTrack(*playingTrack, 0)

		if err = e.DeferCreateMessage(false); err != nil {
			return err
		}

		lyrics, err = lavalyrics.GetCurrentTrackLyrics(ctx, player.Node().Rest(), player.Node().SessionID(), *e.GuildID(), skipTrackSource)
	}
	if err != nil {
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: json.Ptr(fmt.Sprintf("failed to decode track: %s", err)),
		})
		return err
	}

	if lyrics == nil {
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: json.Ptr("no lyrics found"),
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
		Content: json.Ptr(fmt.Sprintf("Loaded lyrics for %s from `%s`(`%s`)", track, lyrics.SourceName, lyrics.Provider)),
		Files: []*discord.File{
			discord.NewFile("lyrics.txt", "", bytes.NewReader([]byte(content))),
		},
	})
	return err
}

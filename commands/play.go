package commands

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/json"
	"github.com/lavalink-devs/lavalink-bot/internal/res"
)

var (
	urlPattern   = regexp.MustCompile("^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")
	queryPattern = regexp.MustCompile(`^(.{2})(search|isrc):(.+)`)
)

func (c *Cmds) Play(e *handler.CommandEvent) error {
	voiceState, ok := c.Client.Caches().VoiceState(*e.GuildID(), e.User().ID)
	if !ok {
		return e.CreateMessage(discord.MessageCreate{
			Content: "You need to be in a voice channel to use this command.",
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	data := e.SlashCommandInteractionData()

	identifier := data.String("identifier")

	if !urlPattern.MatchString(identifier) && !queryPattern.MatchString(identifier) {
		if source, ok := data.OptString("source"); ok {
			identifier = lavalink.SearchType(source).Apply(identifier)
		} else {
			identifier = lavalink.SearchType("dzsearch").Apply(identifier)
		}
	}

	if err := e.DeferCreateMessage(false); err != nil {
		return err
	}
	go loadTracks(c, e, voiceState, identifier)
	return nil
}

func loadTracks(c *Cmds, e *handler.CommandEvent, voiceState discord.VoiceState, identifier string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := c.Lavalink.BestNode().LoadTracks(ctx, identifier)
	if err != nil {
		_, _ = e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: json.Ptr(fmt.Sprintf("Failed to load tracks: %s", err)),
		})
		return
	}

	var (
		tracks       []lavalink.Track
		playlistName string
	)
	switch data := result.Data.(type) {
	case lavalink.Track:
		tracks = append(tracks, data)
	case lavalink.Playlist:
		tracks = append(tracks, data.Tracks...)
		playlistName = data.Info.Name
	case lavalink.Search:
		fmt.Printf("search: %+v\n", data)
		tracks = append(tracks, data[0])
	case lavalink.Empty:
		if _, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: json.Ptr("No matches found"),
		}); err != nil {
			c.Client.Logger().Errorf("failed to update interaction response: %s", err)
		}
		return
	case lavalink.Exception:
		if _, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: json.Ptr(fmt.Sprintf("Failed to load tracks: %s", data.Error())),
		}); err != nil {
			c.Client.Logger().Errorf("failed to update interaction response: %s", err)
		}
		return
	}

	var content string
	if playlistName != "" {
		content = fmt.Sprintf("Loaded playlist **%s** with %d tracks", playlistName, len(tracks))
	} else if len(tracks) == 1 {
		content = fmt.Sprintf("Loaded track **%s**", res.FormatTrack(tracks[0], 0))
	} else {
		content = fmt.Sprintf("Loaded **%d** tracks", len(tracks))
	}

	if _, err = e.UpdateInteractionResponse(discord.MessageUpdate{
		Content: &content,
	}); err != nil {
		c.Client.Logger().Errorf("failed to update interaction response: %s", err)
	}

	if err = c.Client.UpdateVoiceState(context.Background(), *e.GuildID(), voiceState.ChannelID, false, false); err != nil {
		if _, err = e.CreateFollowupMessage(discord.MessageCreate{
			Content: fmt.Sprintf("Failed to join voice channel: %s", err),
		}); err != nil {
			c.Client.Logger().Errorf("failed to create followup message: %s", err)
			return
		}
	}

	player := c.Lavalink.Player(*e.GuildID())
	if player.Track() == nil {
		var track lavalink.Track
		if len(tracks) == 1 {
			track = tracks[0]
			tracks = nil
		} else {
			track, tracks = tracks[0], tracks[1:]
		}

		playCtx, playCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer playCancel()
		if err = player.Update(playCtx, lavalink.WithTrack(track)); err != nil {
			if _, err = e.CreateFollowupMessage(discord.MessageCreate{
				Content: fmt.Sprintf("Failed to play track: %s", err),
			}); err != nil {
				c.Client.Logger().Errorf("failed to create followup message: %s", err)
			}
			return
		}
	}
	if len(tracks) > 0 {
		c.MusicQueue.Add(*e.GuildID(), e.Channel().ID(), tracks...)
	}
}

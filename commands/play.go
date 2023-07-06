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

func (c *Commands) Play(e *handler.CommandEvent) error {
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := c.Lavalink.BestNode().LoadTracks(ctx, identifier)
	if err != nil {
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: json.Ptr(fmt.Sprintf("Failed to load tracks: %s", err)),
		})
		return err
	}

	var (
		tracks         []lavalink.Track
		messageContent string
	)
	switch loadData := result.Data.(type) {
	case lavalink.Track:
		tracks = append(tracks, loadData)
		messageContent = fmt.Sprintf("Loaded track **%s**", res.FormatTrack(loadData, 0))
	case lavalink.Playlist:
		tracks = append(tracks, loadData.Tracks...)
		messageContent = fmt.Sprintf("Loaded playlist **%s** with %d tracks", loadData.Info.Name, len(tracks))
	case lavalink.Search:
		tracks = append(tracks, loadData[0])
		messageContent = fmt.Sprintf("Loaded track **%s** from search", res.FormatTrack(loadData[0], 0))
	case lavalink.Empty:
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: json.Ptr("No matches found"),
		})
		return err
	case lavalink.Exception:
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: json.Ptr(fmt.Sprintf("Failed to load tracks: %s", loadData.Error())),
		})
		return err
	}

	if _, err = e.UpdateInteractionResponse(discord.MessageUpdate{
		Content: &messageContent,
	}); err != nil {
		return err
	}

	if err = c.Client.UpdateVoiceState(context.Background(), *e.GuildID(), voiceState.ChannelID, false, false); err != nil {
		_, err = e.CreateFollowupMessage(discord.MessageCreate{
			Content: fmt.Sprintf("Failed to join voice channel: %s", err),
		})
		return err
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
			_, err = e.CreateFollowupMessage(discord.MessageCreate{
				Content: fmt.Sprintf("Failed to play track: %s", err),
			})
			return err
		}
	}
	if len(tracks) > 0 {
		c.MusicQueue.Add(*e.GuildID(), e.Channel().ID(), tracks...)
	}
	return nil
}

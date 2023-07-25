package commands

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/json"
	"github.com/disgoorg/lavasearch-plugin"
	"github.com/lavalink-devs/lavalink-bot/internal/res"
	"go.deanishe.net/fuzzy"
)

var (
	urlPattern   = regexp.MustCompile("^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")
	queryPattern = regexp.MustCompile(`^(.{2})(search|isrc):(.+)`)
)

func (c *Commands) PlayAutocomplete(e *handler.AutocompleteEvent) error {
	query := e.Data.String("query")
	if query == "" {
		return e.Result(nil)
	}

	source := lavalink.SearchType(e.Data.String("source"))
	if source == "" {
		source = "dzsearch"
	}
	query = source.Apply(query)

	var types []lavasearch.SearchType
	if searchType, ok := e.Data.OptString("type"); ok {
		types = append(types, lavasearch.SearchType(searchType))
	}

	result, err := lavasearch.LoadSearch(c.Lavalink.BestNode().Rest(), query, types)
	if err != nil {
		if errors.Is(err, lavasearch.ErrEmptySearchResult) {
			return e.Result(nil)
		}
		return e.Result([]discord.AutocompleteChoice{
			discord.AutocompleteChoiceString{
				Name:  res.Trim("Failed to load search results: "+err.Error(), 100),
				Value: "error",
			},
		})
	}

	choices := make([]discord.AutocompleteChoice, 0)
	for _, album := range result.Albums {
		if len(choices) >= 5 {
			break
		}
		choices = append(choices, discord.AutocompleteChoiceString{
			Name:  res.Trim("💿 "+album.Name+" - "+album.Artist, 100),
			Value: album.URL,
		})
	}
	for _, artist := range result.Artists {
		if len(choices) >= 10 {
			break
		}
		choices = append(choices, discord.AutocompleteChoiceString{
			Name:  res.Trim("🧑 "+artist.Name, 100),
			Value: artist.URL,
		})
	}
	for _, playlist := range result.Playlists {
		if len(choices) >= 15 {
			break
		}
		choices = append(choices, discord.AutocompleteChoiceString{
			Name:  res.Trim("📜 "+playlist.Name, 100),
			Value: playlist.URL,
		})
	}
	for _, track := range result.Tracks {
		if len(choices) >= 20 {
			break
		}
		choices = append(choices, discord.AutocompleteChoiceString{
			Name:  res.Trim("🎶 "+track.Title+" - "+track.Author, 100),
			Value: track.URI,
		})
	}
	for _, text := range result.Texts {
		if len(choices) >= 25 {
			break
		}
		choices = append(choices, discord.AutocompleteChoiceString{
			Name:  res.Trim("🔎"+text.Text, 100),
			Value: text.Text,
		})
	}

	fuzzy.Sort(Choices(choices), e.Data.String("query"))

	return e.Result(choices)
}

var (
	_ sort.Interface = (*Choices)(nil)
	_ fuzzy.Sortable = (*Choices)(nil)
)

type Choices []discord.AutocompleteChoice

func (c Choices) Keywords(i int) string {
	return string([]rune(c[i].(discord.AutocompleteChoiceString).Name)[2:])
}

func (c Choices) Len() int {
	return len(c)
}

func (c Choices) Less(i, j int) bool {
	return c[i].(discord.AutocompleteChoiceString).Name < c[j].(discord.AutocompleteChoiceString).Name
}

func (c Choices) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c *Commands) Play(e *handler.CommandEvent) error {
	voiceState, ok := c.Client.Caches().VoiceState(*e.GuildID(), e.User().ID)
	if !ok {
		return e.CreateMessage(discord.MessageCreate{
			Content: "You need to be in a voice channel to use this command.",
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	data := e.SlashCommandInteractionData()

	query := data.String("query")

	if !urlPattern.MatchString(query) && !queryPattern.MatchString(query) {
		if source, ok := data.OptString("source"); ok {
			query = lavalink.SearchType(source).Apply(query)
		} else {
			query = lavalink.SearchType("dzsearch").Apply(query)
		}
	}

	if err := e.DeferCreateMessage(false); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := c.Lavalink.BestNode().LoadTracks(ctx, query)
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
		playlistType, playlistName := res.FormatPlaylist(loadData)
		messageContent = fmt.Sprintf("Loaded %s **%s**", playlistType, playlistName)
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

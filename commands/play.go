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
	"github.com/disgoorg/lavasrc-plugin"
	"github.com/disgoorg/snowflake/v2"
	"go.deanishe.net/fuzzy"

	"github.com/lavalink-devs/lavalink-bot/internal/res"
)

var (
	urlPattern   = regexp.MustCompile("^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")
	queryPattern = regexp.MustCompile(`^(.{2})(search|isrc):(.+)`)
)

type UserData struct {
	Requester  snowflake.ID `json:"requester"`
	OriginType string       `json:"origin_type"`
	OriginName string       `json:"origin_name"`
}

func (c *Commands) PlayAutocomplete(e *handler.AutocompleteEvent) error {
	query := e.Data.String("query")
	if query == "" {
		return e.AutocompleteResult(nil)
	}

	if urlPattern.MatchString(query) {
		ctx, cancel := context.WithTimeout(e.Ctx, 10*time.Second)
		defer cancel()
		result, err := c.Lavalink.BestNode().LoadTracks(ctx, query)
		if err != nil {
			return e.AutocompleteResult([]discord.AutocompleteChoice{
				discord.AutocompleteChoiceString{
					Name:  res.Trim("Failed to load track: "+err.Error(), 100),
					Value: "error",
				},
			})
		}

		var name string
		switch loadData := result.Data.(type) {
		case lavalink.Track:
			name = fmt.Sprintf("%s - %s", loadData.Info.Title, loadData.Info.Author)
		case lavalink.Playlist:
			var playlistInfo lavasrc.PlaylistInfo
			_ = loadData.PluginInfo.Unmarshal(&playlistInfo)
			name = fmt.Sprintf("%s: %s - %s", playlistInfo.Type, loadData.Info.Name, playlistInfo.Author)
		case lavalink.Empty:
			return e.AutocompleteResult(nil)
		case lavalink.Exception:
			return e.AutocompleteResult([]discord.AutocompleteChoice{
				discord.AutocompleteChoiceString{
					Name:  res.Trim("Failed to load track: "+loadData.Error(), 100),
					Value: "error",
				},
			})
		}

		return e.AutocompleteResult([]discord.AutocompleteChoice{
			discord.AutocompleteChoiceString{
				Name:  res.Trim("🔗 "+name, 100),
				Value: query,
			},
		})
	}

	source := lavalink.SearchType(e.Data.String("source"))
	if source == "" {
		source = "dzsearch"
	}
	if !e.Data.Bool("raw") {
		query = source.Apply(query)
	}

	var types []lavasearch.SearchType
	if searchType, ok := e.Data.OptString("type"); ok {
		types = append(types, lavasearch.SearchType(searchType))
	}

	ctx, cancel := context.WithTimeout(e.Ctx, 10*time.Second)
	defer cancel()
	result, err := lavasearch.LoadSearch(ctx, c.Lavalink.BestNode().Rest(), query, types)
	if err != nil {
		if errors.Is(err, lavasearch.ErrEmptySearchResult) {
			return e.AutocompleteResult(nil)
		}
		return e.AutocompleteResult([]discord.AutocompleteChoice{
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

		var albumInfo lavasrc.PlaylistInfo
		_ = album.PluginInfo.Unmarshal(&albumInfo)

		choices = append(choices, discord.AutocompleteChoiceString{
			Name:  res.Trim("💿 "+album.Info.Name+" - "+albumInfo.Author, 100),
			Value: albumInfo.URL,
		})
	}
	for _, artist := range result.Artists {
		if len(choices) >= 10 {
			break
		}

		var artistInfo lavasrc.PlaylistInfo
		_ = artist.PluginInfo.Unmarshal(&artistInfo)

		choices = append(choices, discord.AutocompleteChoiceString{
			Name:  res.Trim("🧑 "+artistInfo.Author, 100),
			Value: artistInfo.URL,
		})
	}
	for _, playlist := range result.Playlists {
		if len(choices) >= 15 {
			break
		}

		var playlistInfo lavasrc.PlaylistInfo
		_ = playlist.PluginInfo.Unmarshal(&playlistInfo)

		choices = append(choices, discord.AutocompleteChoiceString{
			Name:  res.Trim("📜 "+playlist.Info.Name, 100),
			Value: playlistInfo.URL,
		})
	}
	for _, track := range result.Tracks {
		if len(choices) >= 20 {
			break
		}

		var trackInfo lavasrc.PlaylistInfo
		_ = track.PluginInfo.Unmarshal(&trackInfo)

		choices = append(choices, discord.AutocompleteChoiceString{
			Name:  res.Trim("🎶 "+track.Info.Title+" - "+track.Info.Author, 100),
			Value: *track.Info.URI,
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

	return e.AutocompleteResult(choices)
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

func (c *Commands) Play(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	voiceState, ok := c.Client.Caches().VoiceState(*e.GuildID(), e.User().ID)
	if !ok {
		return e.CreateMessage(discord.MessageCreate{
			Content: "You need to be in a voice channel to use this command.",
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	query := data.String("query")

	if !data.Bool("raw") {
		if !urlPattern.MatchString(query) && !queryPattern.MatchString(query) {
			if source, ok := data.OptString("source"); ok {
				query = lavalink.SearchType(source).Apply(query)
			} else {
				query = lavalink.SearchType("dzsearch").Apply(query)
			}
		}
	}

	if err := e.DeferCreateMessage(false); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(e.Ctx, 10*time.Second)
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
		userData       = UserData{
			Requester: e.User().ID,
		}
	)
	switch loadData := result.Data.(type) {
	case lavalink.Track:
		tracks = append(tracks, loadData)
		messageContent = fmt.Sprintf("Loaded track **%s**", res.FormatTrack(loadData, 0))
	case lavalink.Playlist:
		tracks = append(tracks, loadData.Tracks...)
		playlistType, playlistName := res.FormatPlaylist(loadData)
		messageContent = fmt.Sprintf("Loaded %s **%s** - `%d tracks`", playlistType, playlistName, len(loadData.Tracks))
		userData.OriginType = playlistType
		userData.OriginName = playlistName
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
			Files:   []*discord.File{res.NewExceptionFile(loadData.CauseStackTrace)},
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

	userDataRaw, _ := json.Marshal(userData)
	for i := range tracks {
		tracks[i].UserData = userDataRaw
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

		c.MusicQueue.Add(*e.GuildID(), e.Channel().ID(), tracks...)

		playCtx, playCancel := context.WithTimeout(e.Ctx, 10*time.Second)
		defer playCancel()
		if err = player.Update(playCtx, lavalink.WithTrack(track)); err != nil {
			_, err = e.CreateFollowupMessage(discord.MessageCreate{
				Content: fmt.Sprintf("Failed to play track: %s", err),
			})
			return err
		}
	} else {
		c.MusicQueue.Add(*e.GuildID(), e.Channel().ID(), tracks...)
	}

	return nil
}

func (c *Commands) PlayTrack(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	voiceState, ok := c.Client.Caches().VoiceState(*e.GuildID(), e.User().ID)
	if !ok {
		return e.CreateMessage(discord.MessageCreate{
			Content: "You need to be in a voice channel to use this command.",
			Flags:   discord.MessageFlagEphemeral,
		})
	}

	encodedTrack := data.String("track")

	if err := e.DeferCreateMessage(false); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(e.Ctx, 10*time.Second)
	defer cancel()
	track, err := c.Lavalink.BestNode().DecodeTrack(ctx, encodedTrack)
	if err != nil {
		_, err = e.UpdateInteractionResponse(discord.MessageUpdate{
			Content: json.Ptr(fmt.Sprintf("Failed to decode track: %s", err)),
		})
		return err
	}

	userDataRaw, _ := json.Marshal(UserData{
		Requester: e.User().ID,
	})
	track.UserData = userDataRaw

	if _, err = e.UpdateInteractionResponse(discord.MessageUpdate{
		Content: json.Ptr(fmt.Sprintf("Loaded track **%s**", res.FormatTrack(*track, 0))),
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
		playCtx, playCancel := context.WithTimeout(e.Ctx, 10*time.Second)
		defer playCancel()
		if err = player.Update(playCtx, lavalink.WithTrack(*track)); err != nil {
			_, err = e.CreateFollowupMessage(discord.MessageCreate{
				Content: fmt.Sprintf("Failed to play track: %s", err),
			})
			return err
		}
	} else {
		c.MusicQueue.Add(*e.GuildID(), e.Channel().ID(), *track)
	}
	return nil
}

package res

import (
	"fmt"

	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/lavasrc-plugin"
)

func FormatTrack(track lavalink.Track, position lavalink.Duration) string {
	var positionStr string
	if position > 0 {
		positionStr = fmt.Sprintf("`%s/%s`", FormatDuration(position), FormatDuration(track.Info.Length))
	} else {
		positionStr = fmt.Sprintf("`%s`", FormatDuration(track.Info.Length))
	}

	if track.Info.URI != nil {
		return fmt.Sprintf("[`%s`](<%s>) - `%s` `%s`", track.Info.Title, *track.Info.URI, track.Info.Author, positionStr)
	}
	return fmt.Sprintf("`%s` - `%s` `%s`", track.Info.Title, track.Info.Author, positionStr)
}

func FormatPlaylist(playlist lavalink.Playlist) (string, string) {
	var lavasrcInfo lavasrc.PlaylistInfo
	_ = playlist.PluginInfo.Unmarshal(&lavasrcInfo)

	playlistType := "playlist"
	if lavasrcInfo.Type != "" {
		playlistType = string(lavasrcInfo.Type)
	}

	name := playlist.Info.Name
	if lavasrcInfo.Author != "" {
		name = lavasrcInfo.Author + " - " + name
	}
	if lavasrcInfo.URL != "" {
		return playlistType, fmt.Sprintf("[`%s`](<%s>) - `%d tracks`", playlist.Info.Name, lavasrcInfo.URL, len(playlist.Tracks))
	}

	return playlistType, fmt.Sprintf("`%s` - `%d tracks`", playlist.Info.Name, len(playlist.Tracks))
}

func FormatDuration(duration lavalink.Duration) string {
	if duration == 0 {
		return "00:00"
	}
	return fmt.Sprintf("%02d:%02d", duration.Minutes(), duration.SecondsPart())
}

func Trim(s string, length int) string {
	r := []rune(s)
	if len(r) > length {
		return string(r[:length-1]) + "…"
	}
	return s
}

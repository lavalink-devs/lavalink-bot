package res

import (
	"fmt"
	"strings"

	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/lavasrc-plugin"
)

func FormatTrack(track lavalink.Track, position lavalink.Duration) string {
	var lavasrcInfo lavasrc.TrackInfo
	_ = track.PluginInfo.Unmarshal(&lavasrcInfo)

	var positionStr string
	if track.Info.IsStream {
		positionStr = "`LIVE`"
	} else if position > 0 {
		positionStr = fmt.Sprintf("`%s/%s`", FormatDuration(position), FormatDuration(track.Info.Length))
	} else {
		positionStr = fmt.Sprintf("`%s`", FormatDuration(track.Info.Length))
	}

	var trackAuthor string
	if track.Info.Author != "Unknown Author" {
		if lavasrcInfo.ArtistURL != "" {
			trackAuthor = fmt.Sprintf("[`%s`](<%s>)", track.Info.Author, lavasrcInfo.ArtistURL)
		} else {
			trackAuthor = fmt.Sprintf("`%s`", track.Info.Author)
		}
	}

	title := track.Info.Title
	title = strings.TrimPrefix(title, "https://")
	title = strings.TrimPrefix(title, "http://")

	trackName := fmt.Sprintf("`%s`", title)
	if track.Info.URI != nil {
		trackName = fmt.Sprintf("[`%s`](<%s>)", title, *track.Info.URI)
	}

	var albumName string
	if lavasrcInfo.AlbumName != "" {
		albumName = fmt.Sprintf("`%s`", lavasrcInfo.AlbumName)
		if lavasrcInfo.AlbumURL != "" {
			albumName = fmt.Sprintf("[`%s`](<%s>)", lavasrcInfo.AlbumName, lavasrcInfo.AlbumURL)
		}
		return fmt.Sprintf("%s - %s %s - %s", trackName, trackAuthor, positionStr, albumName)
	}

	return fmt.Sprintf("%s - %s %s", trackName, trackAuthor, positionStr)
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

package res

import (
	"fmt"

	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/lavasrc-plugin"
)

func FormatPlaylist(playlist lavalink.Playlist) (string, string) {
	var lavasrcInfo lavasrc.PlaylistInfo
	_ = playlist.PluginInfo.Unmarshal(&lavasrcInfo)

	playlistType := "playlist"
	if lavasrcInfo.Type != "" {
		playlistType = string(lavasrcInfo.Type)
	}

	name := playlist.Info.Name
	if lavasrcInfo.Type == lavasrc.PlaylistTypeArtist {
		name = lavasrcInfo.Author
	} else if lavasrcInfo.Author != "" {
		name = lavasrcInfo.Author + " - " + name
	}
	if lavasrcInfo.URL != "" {
		return playlistType, fmt.Sprintf("[`%s`](<%s>) - `%d tracks`", name, lavasrcInfo.URL, len(playlist.Tracks))
	}

	return playlistType, fmt.Sprintf("`%s` - `%d tracks`", name, len(playlist.Tracks))
}

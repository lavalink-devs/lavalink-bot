package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/json"
)

var timeunitChoices = []discord.ApplicationCommandOptionChoiceInt{
	{
		Name:  "Milliseconds",
		Value: int(lavalink.Millisecond),
	},
	{
		Name:  "Seconds",
		Value: int(lavalink.Second),
	},
	{
		Name:  "Minutes",
		Value: int(lavalink.Minute),
	},
	{
		Name:  "Hours",
		Value: int(lavalink.Hour),
	},
	{
		Name:  "Days",
		Value: int(lavalink.Day),
	},
}

var searchSourceChoices = []discord.ApplicationCommandOptionChoiceString{
	{
		Name:  "YouTube",
		Value: string(lavalink.SearchTypeYouTube),
	},
	{
		Name:  "YouTube Music",
		Value: string(lavalink.SearchTypeYouTubeMusic),
	},
	{
		Name:  "Deezer",
		Value: "dzsearch",
	},
	{
		Name:  "Deezer ISRC",
		Value: "dzisrc",
	},
	{
		Name:  "Spotify",
		Value: "spsearch",
	},
	{
		Name:  "AppleMusic",
		Value: "amsearch",
	},
	{
		Name:  "SoundCloud",
		Value: string(lavalink.SearchTypeSoundCloud),
	},
	{
		Name:  "Tidal",
		Value: "tdsearch",
	},
}

var searchTypeChoices = []discord.ApplicationCommandOptionChoiceString{
	{
		Name:  "Track",
		Value: "track",
	},
	{
		Name:  "Album",
		Value: "album",
	},
	{
		Name:  "Artist",
		Value: "artist",
	},
	{
		Name:  "Playlist",
		Value: "playlist",
	},
	{
		Name:  "Text",
		Value: "text",
	},
}

var sponsorblockOptions = []discord.ApplicationCommandOption{
	discord.ApplicationCommandOptionBool{
		Name:        "sponsor",
		Description: "Whether to skip sponsor segments",
	},
	discord.ApplicationCommandOptionBool{
		Name:        "selfpromo",
		Description: "Whether to skip selfpromo segments",
	},
	discord.ApplicationCommandOptionBool{
		Name:        "interaction",
		Description: "Whether to skip interaction segments",
	},
	discord.ApplicationCommandOptionBool{
		Name:        "intro",
		Description: "Whether to skip intro segments",
	},
	discord.ApplicationCommandOptionBool{
		Name:        "outro",
		Description: "Whether to skip outro segments",
	},
	discord.ApplicationCommandOptionBool{
		Name:        "preview",
		Description: "Whether to skip preview segments",
	},
	discord.ApplicationCommandOptionBool{
		Name:        "music_offtopic",
		Description: "Whether to skip music_offtopic segments",
	},
	discord.ApplicationCommandOptionBool{
		Name:        "filler",
		Description: "Whether to skip filler segments",
	},
}

var music = discord.SlashCommandCreate{
	Name:        "music",
	Description: "music commands",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "play",
			Description: "Plays a song from a given identifier",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:         "query",
					Description:  "The query to search or play",
					Required:     true,
					Autocomplete: true,
				},
				discord.ApplicationCommandOptionString{
					Name:        "source",
					Description: "The source to search from",
					Required:    false,
					Choices:     searchSourceChoices,
				},
				discord.ApplicationCommandOptionString{
					Name:        "type",
					Description: "The type of the search",
					Required:    false,
					Choices:     searchTypeChoices,
				},
				discord.ApplicationCommandOptionBool{
					Name:        "raw",
					Description: "Whether to not do any transformation on the query",
					Required:    false,
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "play-track",
			Description: "Plays a given encoded track",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "track",
					Description: "The encoded track to play",
					Required:    true,
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "tts",
			Description: "Crafts a text-to-speech link to play",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "message",
					Description: "The message to play",
					Required:    true,
				},
				discord.ApplicationCommandOptionString{
					Name:        "voice",
					Description: "The voice to use for the message",
					Required:    false,
				},
				discord.ApplicationCommandOptionBool{
					Name:        "translate",
					Description: "Whether to translate the message to english",
					Required:    false,
				},
				discord.ApplicationCommandOptionFloat{
					Name:        "silence",
					Description: "The silence to add before the message",
					Required:    false,
				},
				discord.ApplicationCommandOptionFloat{
					Name:        "speed",
					Description: "The speed of the message",
					Required:    false,
				},
				discord.ApplicationCommandOptionString{
					Name:        "audio-format",
					Description: "The audio format of the message",
					Required:    false,
					Choices: []discord.ApplicationCommandOptionChoiceString{
						{
							Name:  "MP3",
							Value: "mp3",
						},
						{
							Name:  "OGG Opus",
							Value: "ogg_opus",
						},
						{
							Name:  "OGG Vorbis",
							Value: "ogg_vorbis",
						},
						{
							Name:  "WAV",
							Value: "wav",
						},
						{
							Name:  "FLAC",
							Value: "flac",
						},
						{
							Name:  "AAC",
							Value: "aac",
						},
					},
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "stop",
			Description: "Stops the current track",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "disconnect",
			Description: "Disconnects the player",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "skip",
			Description: "Skips the current track",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionInt{
					Name:        "count",
					Description: "The number of tracks to skip",
					Required:    false,
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "pause",
			Description: "Pauses the current track",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "resume",
			Description: "Resumes the current track",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "seek",
			Description: "Seeks to a given position in the current track",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionInt{
					Name:        "position",
					Description: "The position to seek to",
					Required:    true,
				},
				discord.ApplicationCommandOptionInt{
					Name:        "unit",
					Description: "The unit of the position",
					Required:    false,
					Choices:     timeunitChoices,
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "volume",
			Description: "Sets the volume of the player",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionInt{
					Name:        "volume",
					Description: "The volume to set",
					Required:    true,
					MinValue:    json.Ptr(0),
					MaxValue:    json.Ptr(200),
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "shuffle",
			Description: "Shuffles the queue",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "repeat",
			Description: "Sets the repeat mode of the player",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "mode",
					Description: "The repeat mode to set",
					Required:    true,
					Choices: []discord.ApplicationCommandOptionChoiceString{
						{
							Name:  "None",
							Value: "none",
						},
						{
							Name:  "Track",
							Value: "track",
						},
						{
							Name:  "Queue",
							Value: "queue",
						},
					},
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "queue",
			Description: "Shows the current queue",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "now-playing",
			Description: "Shows the current track",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionBool{
					Name:        "raw",
					Description: "Whether to include the raw track & info",
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "lyrics",
			Description: "Shows the lyrics of the current or a given track",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionBool{
					Name:        "track",
					Description: "Whether to include the raw lyrics",
				},
				discord.ApplicationCommandOptionBool{
					Name:        "skip-track-source",
					Description: "Whether to skip the track source to resolve the lyrics",
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "remove",
			Description: "Removes a track from the queue",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionInt{
					Name:         "index",
					Description:  "The index of the track to remove",
					Required:     true,
					Autocomplete: true,
					MinValue:     json.Ptr(0),
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "move",
			Description: "Moves a track in the queue",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionInt{
					Name:         "from",
					Description:  "The index of the track to move",
					Required:     true,
					Autocomplete: true,
				},
				discord.ApplicationCommandOptionInt{
					Name:         "to",
					Description:  "The index to move the track to",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "swap",
			Description: "Swaps two tracks in the queue",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionInt{
					Name:         "first",
					Description:  "The index of the first track to swap",
					Required:     true,
					Autocomplete: true,
				},
				discord.ApplicationCommandOptionInt{
					Name:         "second",
					Description:  "The index of the second track to swap",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "clear",
			Description: "Clears the queue",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "rewind",
			Description: "Rewinds the current track",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionInt{
					Name:        "amount",
					Description: "The amount to rewind",
					Required:    true,
				},
				discord.ApplicationCommandOptionInt{
					Name:        "unit",
					Description: "The unit of the amount",
					Required:    false,
					Choices:     timeunitChoices,
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "forward",
			Description: "Forwards the current track",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionInt{
					Name:        "amount",
					Description: "The amount to forward",
					Required:    true,
				},
				discord.ApplicationCommandOptionInt{
					Name:        "unit",
					Description: "The unit of the amount",
					Required:    false,
					Choices:     timeunitChoices,
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "restart",
			Description: "Restarts the current track",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "effects",
			Description: "Shows the current effects",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "effect",
					Description: "The effect to apply",
					Required:    true,
					Choices: []discord.ApplicationCommandOptionChoiceString{
						{
							Name:  "None",
							Value: string(EffectTypeNone),
						},
						{
							Name:  "Nightcore",
							Value: string(EffectTypeNightcore),
						},
						{
							Name:  "Vaporwave",
							Value: string(EffectTypeVaporwave),
						},
						{
							Name:  "Piano",
							Value: string(EffectTypePiano),
						},
						{
							Name:  "Metal",
							Value: string(EffectTypeMetal),
						},
						{
							Name:  "Bass Boost",
							Value: string(EffectTypeBassBoost),
						},
					},
				},
			},
		},
		discord.ApplicationCommandOptionSubCommandGroup{
			Name:        "sponsorblock",
			Description: "Shows or sets the skipping sponsorblock categories",
			Options: []discord.ApplicationCommandOptionSubCommand{
				{
					Name:        "show",
					Description: "Shows the current skipping sponsorblock categories",
				},
				{
					Name:        "set",
					Description: "Sets the skipping sponsorblock categories",
					Options:     sponsorblockOptions,
				},
			},
		},
	},
}

func (c *Commands) RequirePlayer(next handler.Handler) handler.Handler {
	return func(e *handler.InteractionEvent) error {
		if e.Type() == discord.InteractionTypeApplicationCommand {
			if player := c.Lavalink.ExistingPlayer(*e.GuildID()); player == nil {
				return e.CreateMessage(discord.MessageCreate{
					Content: "No player found",
					Flags:   discord.MessageFlagEphemeral,
				})
			}
		}

		return next(e)
	}
}

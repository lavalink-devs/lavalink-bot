package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/lavalink"
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

var music = discord.SlashCommandCreate{
	Name:        "music",
	Description: "music commands",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "play",
			Description: "Plays a song from a given identifier",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "identifier",
					Description: "The identifier to search or play",
					Required:    true,
				},
				discord.ApplicationCommandOptionString{
					Name:        "source",
					Description: "The source to search from",
					Required:    false,
					Choices: []discord.ApplicationCommandOptionChoiceString{
						{
							Name:  "YouTube",
							Value: string(lavalink.SearchTypeYouTube),
						},
						{
							Name:  "YouTube Music",
							Value: string(lavalink.SearchTypeYouTubeMusic),
						},
						{
							Name:  "SoundCloud",
							Value: string(lavalink.SearchTypeSoundCloud),
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
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "lyrics",
			Description: "Shows the lyrics of the current track",
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
	},
}

func (c *Cmds) RequirePlayer(next handler.Handler) handler.Handler {
	return func(e *events.InteractionCreate) error {
		if e.Type() == discord.InteractionTypeApplicationCommand {
			if player := c.Lavalink.ExistingPlayer(*e.GuildID()); player == nil {
				return e.Respond(discord.InteractionResponseTypeCreateMessage, discord.MessageCreate{
					Content: "No player found",
					Flags:   discord.MessageFlagEphemeral,
				})
			}
		}

		return next(e)
	}
}

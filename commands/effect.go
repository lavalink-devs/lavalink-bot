package commands

import (
	"context"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

type EffectType string

const (
	EffectTypeNone      EffectType = "none"
	EffectTypeNightcore EffectType = "nightcore"
	EffectTypeVaporwave EffectType = "vaporwave"
	EffectTypePiano     EffectType = "piano"
	EffectTypeMetal     EffectType = "metal"
	EffectTypeBassBoost EffectType = "bass-boost"
)

var effects = map[EffectType]lavalink.Filters{
	EffectTypeNone: {},
	EffectTypeNightcore: {
		Equalizer: &lavalink.Equalizer{
			0: -0.075,
			1: 0.125,
			2: 0.125,
		},
		Timescale: &lavalink.Timescale{
			Pitch: 0.95,
			Rate:  1.3,
			Speed: 1,
		},
	},
	EffectTypeVaporwave: {
		Equalizer: &lavalink.Equalizer{
			0: 0.25,
			1: 0.2,
			2: 0.2,
		},
		Timescale: &lavalink.Timescale{
			Pitch: 1,
			Rate:  1,
			Speed: 0.7,
		},
	},
	EffectTypePiano: {
		Equalizer: &lavalink.Equalizer{
			0:  -0.25,
			1:  -0.25,
			2:  -0.125,
			4:  0.25,
			5:  0.25,
			7:  -0.25,
			8:  -0.25,
			11: 0.5,
			12: 0.25,
			13: -0.025,
		},
	},
	EffectTypeMetal: {
		Equalizer: &lavalink.Equalizer{
			1:  0.1,
			2:  0.1,
			3:  0.15,
			4:  0.13,
			5:  0.1,
			7:  0.125,
			8:  0.175,
			9:  0.175,
			10: 0.125,
			11: 0.125,
			12: 0.1,
			13: 0.075,
		},
	},
	EffectTypeBassBoost: {
		Equalizer: &lavalink.Equalizer{
			0:  0.2,
			1:  0.15,
			2:  0.1,
			3:  0.05,
			4:  0.0,
			5:  -0.05,
			6:  -0.1,
			7:  -0.1,
			8:  -0.1,
			9:  -0.1,
			10: -0.1,
			11: -0.1,
			12: -0.1,
			13: -0.1,
			14: -0.1,
		},
	},
}

func (c *Commands) Effects(e *handler.CommandEvent) error {
	effectType := EffectType(e.SlashCommandInteractionData().String("effect"))
	if err := c.Lavalink.ExistingPlayer(*e.GuildID()).Update(context.Background(), lavalink.WithFilters(effects[effectType])); err != nil {
		return e.CreateMessage(discord.MessageCreate{
			Content: fmt.Sprintf("failed to appply effect: `%s`", err),
			Flags:   discord.MessageFlagEphemeral,
		})
	}
	return e.CreateMessage(discord.MessageCreate{
		Content: fmt.Sprintf("applied effect: `%s`", effectType),
	})
}

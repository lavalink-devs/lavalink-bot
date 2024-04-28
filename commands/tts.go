package commands

import (
	"net/url"
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func (c *Commands) TTS(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	message := data.String("message")

	tts := url.URL{
		Scheme: "ftts",
		Host:   message,
	}

	query := url.Values{}
	if voice, ok := data.OptString("voice"); ok {
		query.Add("voice", voice)
	}
	if translate, ok := data.OptBool("translate"); ok {
		query.Add("translate", strconv.FormatBool(translate))
	}
	if silence, ok := data.OptFloat("silence"); ok {
		query.Add("silence", strconv.FormatFloat(silence, 'f', -1, 64))
	}
	if speed, ok := data.OptFloat("speed"); ok {
		query.Add("speed", strconv.FormatFloat(speed, 'f', -1, 64))
	}
	if audioFormat, ok := data.OptString("audio-format"); ok {
		query.Add("audio_format", audioFormat)
	}

	tts.RawQuery = query.Encode()

	return e.CreateMessage(discord.MessageCreate{
		Content: tts.String(),
	})
}

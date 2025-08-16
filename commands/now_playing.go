package commands

import (
	"bytes"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"

	"github.com/lavalink-devs/lavalink-bot/internal/res"
)

func (c *Commands) NowPlaying(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	player := c.Lavalink.Player(*e.GuildID())
	track := player.Track()
	if track == nil {
		return e.CreateMessage(discord.MessageCreate{
			Content: "There is no track playing.",
			Flags:   discord.MessageFlagEphemeral,
		})
	}
	content := fmt.Sprintf("Now playing: %s", res.FormatTrack(*track, player.Position()))
	var userData UserData
	_ = track.UserData.Unmarshal(&userData)
	if userData.Requester > 0 {
		content += "\nRequested by: " + discord.UserMention(userData.Requester)
	}
	if userData.OriginType == "playlist" {
		content += fmt.Sprintf("\nFrom: %s", userData.OriginName)
	}

	var files []*discord.File
	if data.Bool("raw") {
		decodedData := MarshalNoEscape(track)
		files = append(files, &discord.File{
			Name:   "track.json",
			Reader: bytes.NewReader(decodedData),
		})
	}

	return e.CreateMessage(discord.MessageCreate{
		Content:         content,
		AllowedMentions: &discord.AllowedMentions{},
		Files:           files,
	})
}

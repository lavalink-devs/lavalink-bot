package handlers

import (
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/json"
	"github.com/disgoorg/lavalyrics-plugin"
	"github.com/topi314/tint"
)

func (h *Handlers) OnLyricsFound(p disgolink.Player, event lavalyrics.LyricsFoundEvent) {
	channelID := h.MusicQueue.ChannelID(p.GuildID())
	if channelID == 0 {
		return
	}

	lyricsMessageID := h.MusicQueue.LyricsMessageID(p.GuildID())
	if lyricsMessageID == 0 {
		msg, err := h.Client.Rest().CreateMessage(channelID, discord.MessageCreate{
			Content: "Lyrics found",
		})
		if err != nil {
			slog.Error("failed to send message", tint.Err(err))
			return
		}
		h.MusicQueue.SetLyricsMessageID(p.GuildID(), msg.ID)
		return
	}

	if _, err := h.Client.Rest().UpdateMessage(channelID, lyricsMessageID, discord.MessageUpdate{
		Content: json.Ptr("Lyrics found"),
	}); err != nil {
		slog.Error("failed to update message", tint.Err(err))
	}
}

func (h *Handlers) OnLyricsNotFound(p disgolink.Player, event lavalyrics.LyricsNotFoundEvent) {
	channelID := h.MusicQueue.ChannelID(p.GuildID())
	if channelID == 0 {
		return
	}

	lyricsMessageID := h.MusicQueue.LyricsMessageID(p.GuildID())
	if lyricsMessageID == 0 {
		msg, err := h.Client.Rest().CreateMessage(channelID, discord.MessageCreate{
			Content: "Lyrics not found",
		})
		if err != nil {
			slog.Error("failed to send message", tint.Err(err))
			return
		}
		h.MusicQueue.SetLyricsMessageID(p.GuildID(), msg.ID)
		return
	}

	if _, err := h.Client.Rest().UpdateMessage(channelID, lyricsMessageID, discord.MessageUpdate{
		Content: json.Ptr("Lyrics not found"),
	}); err != nil {
		slog.Error("failed to update message", tint.Err(err))
	}
}

func (h *Handlers) OnLyricsLine(p disgolink.Player, event lavalyrics.LyricsLineEvent) {
	channelID := h.MusicQueue.ChannelID(p.GuildID())
	if channelID == 0 {
		return
	}

	lyricsMessageID := h.MusicQueue.LyricsMessageID(p.GuildID())
	if lyricsMessageID == 0 {
		return
	}

	if _, err := h.Client.Rest().UpdateMessage(channelID, lyricsMessageID, discord.MessageUpdate{
		Content: json.Ptr(fmt.Sprintf("Line(`%s`): %s", event.Line.Timestamp, event.Line.Line)),
	}); err != nil {
		slog.Error("failed to update message", tint.Err(err))
	}
}

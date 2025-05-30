package handlers

import (
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/json"
	"github.com/disgoorg/lavalyrics-plugin"
	"github.com/topi314/tint"

	"github.com/lavalink-devs/lavalink-bot/internal/res"
	"github.com/lavalink-devs/lavalink-bot/lavalinkbot"
)

func (h *Handlers) OnLyricsFound(p disgolink.Player, event lavalyrics.LyricsFoundEvent) {
	channelID := h.MusicQueue.ChannelID(p.GuildID())
	if channelID == 0 {
		return
	}

	track := p.Track()
	if track == nil {
		return
	}

	content := fmt.Sprintf("Lyrics found for %s from `%s`(`%s`)", res.FormatTrack(*track, 0), event.Lyrics.SourceName, event.Lyrics.Provider)

	msg, err := h.Client.Rest().CreateMessage(channelID, discord.MessageCreate{
		Content: content,
	})
	if err != nil {
		slog.Error("failed to send message", tint.Err(err))
		return
	}
	h.MusicQueue.SetLyrics(p.GuildID(), lavalinkbot.Lyrics{
		MessageID:   msg.ID,
		BaseMessage: content,
	})
}

func (h *Handlers) OnLyricsNotFound(p disgolink.Player, event lavalyrics.LyricsNotFoundEvent) {
	channelID := h.MusicQueue.ChannelID(p.GuildID())
	if channelID == 0 {
		return
	}

	track := p.Track()
	if track == nil {
		return
	}

	content := fmt.Sprintf("Lyrics not found for %s", res.FormatTrack(*track, 0))

	_, err := h.Client.Rest().CreateMessage(channelID, discord.MessageCreate{
		Content: content,
	})
	if err != nil {
		slog.Error("failed to send message", tint.Err(err))
		return
	}
}

func (h *Handlers) OnLyricsLine(p disgolink.Player, event lavalyrics.LyricsLineEvent) {
	channelID := h.MusicQueue.ChannelID(p.GuildID())
	if channelID == 0 {
		return
	}

	lyrics := h.MusicQueue.Lyrics(p.GuildID())
	if lyrics == nil {
		return
	}
	content := fmt.Sprintf("%s\nLine(i: `%d`, t: `%s`, s: `%t`): %s", lyrics.BaseMessage, event.LineIndex, event.Line.Timestamp, event.Skipped, event.Line.Line)

	if _, err := h.Client.Rest().UpdateMessage(channelID, lyrics.MessageID, discord.MessageUpdate{
		Content: json.Ptr(content),
	}); err != nil {
		slog.Error("failed to update message", tint.Err(err))
	}
}

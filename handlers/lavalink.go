package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/sponsorblock-plugin"
	"github.com/lavalink-devs/lavalink-bot/internal/res"
	"github.com/topi314/tint"
)

func (h *Handlers) OnVoiceStateUpdate(event *events.GuildVoiceStateUpdate) {
	if event.VoiceState.UserID != h.Client.ApplicationID() {
		_, ok := h.Client.Caches().VoiceState(event.VoiceState.GuildID, h.Client.ApplicationID())
		if !ok || event.OldVoiceState.ChannelID == nil {
			return
		}
		var voiceStates int
		h.Client.Caches().VoiceStatesForEach(event.VoiceState.GuildID, func(vs discord.VoiceState) {
			if *vs.ChannelID == *event.OldVoiceState.ChannelID {
				voiceStates++
			}
		})
		if voiceStates <= 1 {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := h.Client.UpdateVoiceState(ctx, event.VoiceState.GuildID, nil, false, false); err != nil {
				slog.Error("failed to disconnect from voice channel", tint.Err(err))
			}
		}
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	h.Lavalink.OnVoiceStateUpdate(ctx, event.VoiceState.GuildID, event.VoiceState.ChannelID, event.VoiceState.SessionID)
	if event.VoiceState.ChannelID == nil {
		h.MusicQueue.Delete(event.VoiceState.GuildID)
	}
}

func (h *Handlers) OnVoiceServerUpdate(event *events.VoiceServerUpdate) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	h.Lavalink.OnVoiceServerUpdate(ctx, event.GuildID, event.Token, *event.Endpoint)
}

func (h *Handlers) OnTrackStart(p disgolink.Player, event lavalink.TrackStartEvent) {
	channelID := h.MusicQueue.ChannelID(p.GuildID())
	if channelID == 0 {
		return
	}
	if _, err := h.Client.Rest().CreateMessage(channelID, discord.MessageCreate{
		Content: "Now playing: " + res.FormatTrack(event.Track, 0),
	}); err != nil {
		slog.Error("failed to send message", tint.Err(err))
	}
}

func (h *Handlers) OnTrackEnd(p disgolink.Player, event lavalink.TrackEndEvent) {
	if !event.Reason.MayStartNext() {
		return
	}
	track, ok := h.MusicQueue.Next(p.GuildID())
	if !ok {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := p.Update(ctx, lavalink.WithTrack(track)); err != nil {
		channelID := h.MusicQueue.ChannelID(p.GuildID())
		if channelID == 0 {
			return
		}
		if _, err = h.Client.Rest().CreateMessage(channelID, discord.MessageCreate{
			Content: "failed to start next track: " + err.Error(),
		}); err != nil {
			slog.Error("failed to send message", tint.Err(err))
		}
	}
}

func (h *Handlers) OnTrackException(p disgolink.Player, event lavalink.TrackExceptionEvent) {
	channelID := h.MusicQueue.ChannelID(p.GuildID())
	if channelID == 0 {
		return
	}
	if _, err := h.Client.Rest().CreateMessage(channelID, discord.MessageCreate{
		Content: "Track exception: " + event.Exception.Error(),
	}); err != nil {
		slog.Error("failed to send message", tint.Err(err))
	}
}

func (h *Handlers) OnTrackStuck(p disgolink.Player, event lavalink.TrackStuckEvent) {
	channelID := h.MusicQueue.ChannelID(p.GuildID())
	if channelID == 0 {
		return
	}
	if _, err := h.Client.Rest().CreateMessage(channelID, discord.MessageCreate{
		Content: "Track stuck: " + event.Track.Info.Title,
	}); err != nil {
		slog.Error("failed to send message", tint.Err(err))
	}
}

func (h *Handlers) OnWebSocketClosed(p disgolink.Player, event lavalink.WebSocketClosedEvent) {
	slog.Info("websocket closed: %s, code: %d, reason: %s", event.GuildID(), event.Code, event.Reason)
	slog.Info("websocket closed", slog.Int64("guild_id", int64(event.GuildID())), slog.Int("code", event.Code), slog.String("reason", event.Reason))
}

func (h *Handlers) OnUnknownEvent(p disgolink.Player, event lavalink.UnknownEvent) {
	slog.Info("unknown event", slog.String("event", string(event.Type())), slog.Int64("guild_id", int64(event.GuildID())), slog.String("data", string(event.Data)))
}

func (h *Handlers) OnUnknownMessage(p disgolink.Player, event lavalink.UnknownMessage) {
	slog.Info("unknown message", slog.String("op", string(event.Op())), slog.String("data", string(event.Data)))
}

func (h *Handlers) OnSegmentsLoaded(p disgolink.Player, event sponsorblock.SegmentsLoadedEvent) {
	channelID := h.MusicQueue.ChannelID(p.GuildID())
	if channelID == 0 {
		return
	}

	content := "Segments loaded:\n"
	for i, segment := range event.Segments {
		line := fmt.Sprintf("%d. %s: %s - %s\n", i+1, segment.Category, res.FormatDuration(segment.Start), res.FormatDuration(segment.End))
		if len(content)+len(line) > 2000 {
			content += "..."
			break
		}
		content += line
	}
	if _, err := h.Client.Rest().CreateMessage(channelID, discord.MessageCreate{
		Content: content,
	}); err != nil {
		slog.Error("failed to send message", tint.Err(err))
	}
}

func (h *Handlers) OndSegmentSkipped(p disgolink.Player, event sponsorblock.SegmentSkippedEvent) {
	channelID := h.MusicQueue.ChannelID(p.GuildID())
	if channelID == 0 {
		return
	}
	if _, err := h.Client.Rest().CreateMessage(channelID, discord.MessageCreate{
		Content: fmt.Sprintf("Segment skipped: %s: %s - %s", event.Segment.Category, res.FormatDuration(event.Segment.Start), res.FormatDuration(event.Segment.End)),
	}); err != nil {
		slog.Error("failed to send message", tint.Err(err))
	}
}

func (h *Handlers) OnChaptersLoaded(p disgolink.Player, event sponsorblock.ChaptersLoadedEvent) {
	channelID := h.MusicQueue.ChannelID(p.GuildID())
	if channelID == 0 {
		return
	}

	content := "Chapters loaded:\n"
	for i, chapter := range event.Chapters {
		line := fmt.Sprintf("%d. %s: %s - %s\n", i+1, chapter.Name, res.FormatDuration(chapter.Start), res.FormatDuration(chapter.End))
		if len(content)+len(line) > 2000 {
			content += "..."
			break
		}
		content += line
	}
	if _, err := h.Client.Rest().CreateMessage(channelID, discord.MessageCreate{
		Content: content,
	}); err != nil {
		slog.Error("failed to send message", tint.Err(err))
	}
}

func (h *Handlers) OnChapterStarted(p disgolink.Player, event sponsorblock.ChapterStartedEvent) {
	channelID := h.MusicQueue.ChannelID(p.GuildID())
	if channelID == 0 {
		return
	}
	if _, err := h.Client.Rest().CreateMessage(channelID, discord.MessageCreate{
		Content: fmt.Sprintf("Chapter started: %s: %s - %s", event.Chapter.Name, res.FormatDuration(event.Chapter.Start), res.FormatDuration(event.Chapter.End)),
	}); err != nil {
		slog.Error("failed to send message", tint.Err(err))
	}
}

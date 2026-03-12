package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/v4/disgolink"
	"github.com/disgoorg/sponsorblock-plugin"

	"github.com/lavalink-devs/lavalink-bot/commands"
	"github.com/lavalink-devs/lavalink-bot/internal/res"
)

func (h *Handlers) OnVoiceStateUpdate(event *events.GuildVoiceStateUpdate) {
	if event.VoiceState.UserID != h.Client.ApplicationID {
		_, ok := h.Client.Caches.VoiceState(event.VoiceState.GuildID, h.Client.ApplicationID)
		if !ok || event.OldVoiceState.ChannelID == nil {
			return
		}
		var voiceStates int
		for vs := range h.Client.Caches.VoiceStates(event.VoiceState.GuildID) {
			if *vs.ChannelID == *event.OldVoiceState.ChannelID {
				voiceStates++
			}
		}
		if voiceStates <= 1 {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := h.Client.UpdateVoiceState(ctx, event.VoiceState.GuildID, nil, false, false); err != nil {
				slog.Error("failed to disconnect from voice channel", slog.Any("err", err))
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
	if event.Endpoint == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	h.Lavalink.OnVoiceServerUpdate(ctx, event.GuildID, event.Token, *event.Endpoint)
}

func (h *Handlers) OnTrackStart(event *disgolink.PlayerTrackStartEvent) {
	channelID := h.MusicQueue.ChannelID(event.GuildID)
	if channelID == 0 {
		return
	}

	content := "Now playing: " + res.FormatTrack(event.Track, 0)
	var userData commands.UserData
	_ = event.Track.UserData.Unmarshal(&userData)
	if userData.Requester > 0 {
		content += "\nRequested by: " + discord.UserMention(userData.Requester)
	}
	if userData.OriginType == "playlist" {
		content += fmt.Sprintf("\nFrom: %s", userData.OriginName)
	}

	if _, err := h.Client.Rest.CreateMessage(channelID, discord.MessageCreate{
		Content:         content,
		AllowedMentions: &discord.AllowedMentions{},
	}); err != nil {
		slog.Error("failed to send message", slog.Any("err", err))
	}
}

func (h *Handlers) OnTrackEnd(event *disgolink.PlayerTrackEndEvent) {
	if !event.Reason.MayStartNext() {
		return
	}
	track, ok := h.MusicQueue.Next(event.GuildID)
	if !ok {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := event.Player.Update(ctx, disgolink.WithTrack(track)); err != nil {
		channelID := h.MusicQueue.ChannelID(event.GuildID)
		if channelID == 0 {
			return
		}
		if _, err = h.Client.Rest.CreateMessage(channelID, discord.MessageCreate{
			Content:         "failed to start next track: " + err.Error(),
			AllowedMentions: &discord.AllowedMentions{},
		}); err != nil {
			slog.Error("failed to send message", slog.Any("err", err))
		}
	}
}

func (h *Handlers) OnTrackException(event *disgolink.PlayerTrackExceptionEvent) {
	channelID := h.MusicQueue.ChannelID(event.GuildID)
	if channelID == 0 {
		return
	}
	if _, err := h.Client.Rest.CreateMessage(channelID, discord.MessageCreate{
		Content:         "Track exception: " + event.Exception.Error(),
		Files:           []*discord.File{res.NewExceptionFile(event.Exception.CauseStackTrace)},
		AllowedMentions: &discord.AllowedMentions{},
	}); err != nil {
		slog.Error("failed to send message", slog.Any("err", err))
	}
}

func (h *Handlers) OnTrackStuck(event *disgolink.PlayerTrackStuckEvent) {
	channelID := h.MusicQueue.ChannelID(event.GuildID)
	if channelID == 0 {
		return
	}
	if _, err := h.Client.Rest.CreateMessage(channelID, discord.MessageCreate{
		Content:         "Track stuck: " + event.Track.Info.Title,
		AllowedMentions: &discord.AllowedMentions{},
	}); err != nil {
		slog.Error("failed to send message", slog.Any("err", err))
	}
}

func (h *Handlers) OnWebSocketClosed(event *disgolink.PlayerWebSocketClosedEvent) {
	slog.Info("websocket closed", slog.Int64("guild_id", int64(event.GuildID)), slog.Int("code", event.Code), slog.String("reason", event.Reason))
}

func (h *Handlers) OnUnknownPlayerEvent(event *disgolink.UnknownPlayerEvent) {
	slog.Info("unknown event", slog.String("event", string(event.EventType)), slog.Int64("guild_id", int64(event.GuildID)), slog.String("data", string(event.Data)))
}

func (h *Handlers) OnUnknownEvent(event *disgolink.UnknownEvent) {
	slog.Info("unknown message", slog.String("op", string(event.Op())), slog.String("data", string(event.Data)))
}

func (h *Handlers) OnSegmentsLoaded(event *sponsorblock.SegmentsLoadedEvent) {
	channelID := h.MusicQueue.ChannelID(event.GuildID)
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
	if _, err := h.Client.Rest.CreateMessage(channelID, discord.MessageCreate{
		Content:         content,
		AllowedMentions: &discord.AllowedMentions{},
	}); err != nil {
		slog.Error("failed to send message", slog.Any("err", err))
	}
}

func (h *Handlers) OndSegmentSkipped(event *sponsorblock.SegmentSkippedEvent) {
	channelID := h.MusicQueue.ChannelID(event.GuildID)
	if channelID == 0 {
		return
	}
	if _, err := h.Client.Rest.CreateMessage(channelID, discord.MessageCreate{
		Content:         fmt.Sprintf("Segment skipped: %s: %s - %s", event.Segment.Category, res.FormatDuration(event.Segment.Start), res.FormatDuration(event.Segment.End)),
		AllowedMentions: &discord.AllowedMentions{},
	}); err != nil {
		slog.Error("failed to send message", slog.Any("err", err))
	}
}

func (h *Handlers) OnChaptersLoaded(event *sponsorblock.ChaptersLoadedEvent) {
	channelID := h.MusicQueue.ChannelID(event.GuildID)
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
	if _, err := h.Client.Rest.CreateMessage(channelID, discord.MessageCreate{
		Content:         content,
		AllowedMentions: &discord.AllowedMentions{},
	}); err != nil {
		slog.Error("failed to send message", slog.Any("err", err))
	}
}

func (h *Handlers) OnChapterStarted(event *sponsorblock.ChapterStartedEvent) {
	channelID := h.MusicQueue.ChannelID(event.GuildID)
	if channelID == 0 {
		return
	}
	if _, err := h.Client.Rest.CreateMessage(channelID, discord.MessageCreate{
		Content:         fmt.Sprintf("Chapter started: %s: %s - %s", event.Chapter.Name, res.FormatDuration(event.Chapter.Start), res.FormatDuration(event.Chapter.End)),
		AllowedMentions: &discord.AllowedMentions{},
	}); err != nil {
		slog.Error("failed to send message", slog.Any("err", err))
	}
}

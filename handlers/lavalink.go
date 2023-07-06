package handlers

import (
	"context"
	"github.com/disgoorg/disgo/events"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/log"
	"github.com/lavalink-devs/lavalink-bot/internal/res"
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
				log.Errorf("failed to disconnect from voice channel: %s", err)
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
		h.Client.Logger().Error("failed to send message: ", err)
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
			h.Client.Logger().Error("failed to send message: ", err)
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
		h.Client.Logger().Error("failed to send message: ", err)
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
		h.Client.Logger().Error("failed to send message: ", err)
	}
}

func (h *Handlers) OnWebSocketClosed(p disgolink.Player, event lavalink.WebSocketClosedEvent) {
	channelID := h.MusicQueue.ChannelID(p.GuildID())
	if channelID == 0 {
		return
	}
	if _, err := h.Client.Rest().CreateMessage(channelID, discord.MessageCreate{
		Content: "Websocket closed: " + event.Reason,
	}); err != nil {
		h.Client.Logger().Error("failed to send message: ", err)
	}
}

func (h *Handlers) OnUnknownEvent(p disgolink.Player, event lavalink.UnknownEvent) {
	log.Infof("unknown event: %s, guild_id: %s, data: %s", event.Type(), event.GuildID(), string(event.Data))
}

func (h *Handlers) OnUnknownMessage(p disgolink.Player, event lavalink.UnknownMessage) {
	log.Infof("unknown message: %s, data: %s", event.Op(), string(event.Data))
}

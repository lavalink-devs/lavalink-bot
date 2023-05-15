package handlers

import (
	"context"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/log"
	"github.com/lavalink-devs/lavalink-bot/internal/res"
)

func (h *Hdlr) OnTrackStart(p disgolink.Player, event lavalink.TrackStartEvent) {
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

func (h *Hdlr) OnTrackEnd(p disgolink.Player, event lavalink.TrackEndEvent) {
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

func (h *Hdlr) OnTrackException(p disgolink.Player, event lavalink.TrackExceptionEvent) {
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

func (h *Hdlr) OnTrackStuck(p disgolink.Player, event lavalink.TrackStuckEvent) {
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

func (h *Hdlr) OnWebSocketClosed(p disgolink.Player, event lavalink.WebSocketClosedEvent) {
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

func (h *Hdlr) OnUnknownEvent(p disgolink.Player, event lavalink.UnknownEvent) {
	log.Infof("unknown event: %s, guild_id: %s, data: %s", event.Type(), event.GuildID(), string(event.Data))
}

func (h *Hdlr) OnUnknownMessage(p disgolink.Player, event lavalink.UnknownMessage) {
	log.Infof("unknown message: %s, data: %s", event.Op(), string(event.Data))
}

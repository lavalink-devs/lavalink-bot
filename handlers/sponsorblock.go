package handlers

import (
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/sponsorblock-plugin"
	"github.com/lavalink-devs/lavalink-bot/internal/res"
	"github.com/topi314/tint"
)

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
		Content:         content,
		AllowedMentions: &discord.AllowedMentions{},
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
		Content:         fmt.Sprintf("Segment skipped: %s: %s - %s", event.Segment.Category, res.FormatDuration(event.Segment.Start), res.FormatDuration(event.Segment.End)),
		AllowedMentions: &discord.AllowedMentions{},
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
		Content:         content,
		AllowedMentions: &discord.AllowedMentions{},
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
		Content:         fmt.Sprintf("Chapter started: %s: %s - %s", event.Chapter.Name, res.FormatDuration(event.Chapter.Start), res.FormatDuration(event.Chapter.End)),
		AllowedMentions: &discord.AllowedMentions{},
	}); err != nil {
		slog.Error("failed to send message", tint.Err(err))
	}
}


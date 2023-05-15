package lavalinkbot

import (
	"context"
	"net/http"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/log"
	"github.com/google/go-github/v52/github"
)

type Bot struct {
	Cfg        Config
	Client     bot.Client
	Lavalink   disgolink.Client
	GitHub     *github.Client
	HTTPClient *http.Client
	MusicQueue *PlayerManager
}

func (b *Bot) Start() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := b.Client.OpenGateway(ctx); err != nil {
		return err
	}
	for _, node := range b.Cfg.Nodes {
		if _, err := b.Lavalink.AddNode(ctx, node.ToNodeConfig()); err != nil {
			log.Errorf("failed to add lavalink node %s: %s", node.Name, err)
		} else {
			log.Infof("added lavalink node: %s", node.Name)
		}
	}

	return nil
}

func (b *Bot) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	b.HTTPClient.CloseIdleConnections()
	b.Lavalink.Close()
	b.Client.Close(ctx)
}

func (b *Bot) OnVoiceStateUpdate(event *events.GuildVoiceStateUpdate) {
	if event.VoiceState.UserID != b.Client.ApplicationID() {
		return
	}
	b.Lavalink.OnVoiceStateUpdate(context.TODO(), event.VoiceState.GuildID, event.VoiceState.ChannelID, event.VoiceState.SessionID)
	if event.VoiceState.ChannelID == nil {
		b.MusicQueue.Delete(event.VoiceState.GuildID)
	}
}

func (b *Bot) OnVoiceServerUpdate(event *events.VoiceServerUpdate) {
	b.Lavalink.OnVoiceServerUpdate(context.TODO(), event.GuildID, event.Token, *event.Endpoint)
}

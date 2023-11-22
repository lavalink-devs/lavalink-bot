package lavalinkbot

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/log"
	"github.com/google/go-github/v52/github"

	"github.com/lavalink-devs/lavalink-bot/internal/maven"
)

type Bot struct {
	Cfg        Config
	Client     bot.Client
	Maven      *maven.Maven
	Lavalink   disgolink.Client
	GitHub     *github.Client
	MusicQueue *PlayerManager
}

func (b *Bot) Start() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := b.Client.OpenGateway(ctx); err != nil {
		return err
	}

	var wg sync.WaitGroup
	for i := range b.Cfg.Nodes {
		wg.Add(1)
		go func(node NodeConfig) {
			defer wg.Done()
			if _, err := b.Lavalink.AddNode(ctx, node.ToNodeConfig()); err != nil {
				log.Errorf("failed to add lavalink node %s: %s", node.Name, err)
			} else {
				log.Infof("added lavalink node: %s", node.Name)
			}
		}(b.Cfg.Nodes[i])
	}

	wg.Wait()
	if node := b.Lavalink.BestNode(); node == nil {
		log.Error("no node connected")
		os.Exit(-1)
	}

	return nil
}

func (b *Bot) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	b.Maven.Close()
	b.Lavalink.Close()
	b.Client.Close(ctx)
}

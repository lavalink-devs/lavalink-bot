package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/disgoorg/disgo/handler/middleware"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/log"
	"github.com/google/go-github/v52/github"
	"github.com/lavalink-devs/lavalink-bot/commands"
	"github.com/lavalink-devs/lavalink-bot/handlers"
	"github.com/lavalink-devs/lavalink-bot/lavalinkbot"
)

func main() {
	path := flag.String("config", "config.yml", "path to config.yml")
	flag.Parse()

	cfg, err := lavalinkbot.ReadConfig(*path)
	if err != nil {
		log.Fatal("failed to read config: ", err)
	}

	log.SetFlags(cfg.Log.Flags())
	log.SetLevel(cfg.Log.Level)
	log.Info("starting lavalink-bot...")
	log.Info("disgo version: ", disgo.Version)
	log.Info("disgolink version: ", disgolink.Version)
	log.Info("loading config from: ", *path)

	b := &lavalinkbot.Bot{
		Cfg:    cfg,
		GitHub: github.NewClient(nil),
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		MusicQueue: lavalinkbot.NewPlayerManager(),
	}

	cmds := &commands.Commands{Bot: b}
	r := handler.New()
	r.Use(middleware.Go)
	r.Command("/info", cmds.Info)
	r.Command("/latest", cmds.Latest)
	r.Autocomplete("/latest", cmds.LatestAutocomplete)
	r.Command("/decode", cmds.Decode)

	r.Route("/music", func(r handler.Router) {
		r.Command("/play", cmds.Play)
		r.Autocomplete("/play", cmds.PlayAutocomplete)
		r.Group(func(r handler.Router) {
			r.Use(cmds.RequirePlayer)

			r.Command("/stop", cmds.Stop)
			r.Command("/disconnect", cmds.Disconnect)
			r.Command("/skip", cmds.Skip)
			r.Command("/pause", cmds.Pause)
			r.Command("/resume", cmds.Resume)
			r.Command("/seek", cmds.Seek)
			// r.Command("/volume", cmds.Volume)
			// r.Command("/shuffle", cmds.Shuffle)
			// r.Command("/repeat", cmds.Repeat)
			r.Command("/queue", cmds.Queue)
			r.Command("/now-playing", cmds.NowPlaying)
			// r.Command("/lyrics", cmds.Lyrics)
			// r.Command("/remove", cmds.Remove)
			// r.Command("/move", cmds.Move)
			// r.Command("/swap", cmds.Swap)
			// r.Command("/clear", cmds.Clear)
			// r.Command("/rewind", cmds.Rewind)
			// r.Command("/forward", cmds.Forward)
			// r.Command("/restart", cmds.Restart)
			r.Command("/effects", cmds.Effects)
		})
	})

	hdlr := &handlers.Handlers{Bot: b}

	if b.Client, err = disgo.New(cfg.Bot.Token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentGuilds, gateway.IntentGuildVoiceStates),
		),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagVoiceStates),
		),
		bot.WithEventListeners(r),
		bot.WithEventListenerFunc(hdlr.OnVoiceStateUpdate),
		bot.WithEventListenerFunc(hdlr.OnVoiceServerUpdate),
	); err != nil {
		log.Fatal("failed to create disgo client: ", err)
	}

	if err = handler.SyncCommands(b.Client, commands.CommandCreates, b.Cfg.Bot.GuildIDs); err != nil {
		log.Errorf("failed to sync commands: %s", err)
	}

	if b.Lavalink = disgolink.New(b.Client.ApplicationID(),
		disgolink.WithListenerFunc(hdlr.OnTrackStart),
		disgolink.WithListenerFunc(hdlr.OnTrackEnd),
		disgolink.WithListenerFunc(hdlr.OnTrackException),
		disgolink.WithListenerFunc(hdlr.OnTrackStuck),
		disgolink.WithListenerFunc(hdlr.OnWebSocketClosed),
		disgolink.WithListenerFunc(hdlr.OnUnknownEvent),
	); err != nil {
		log.Fatal("failed to create disgolink client: ", err)
	}

	if err = b.Start(); err != nil {
		log.Fatal("failed to start bot: ", err)
	}
	defer b.Stop()

	log.Info("lavalink-bot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

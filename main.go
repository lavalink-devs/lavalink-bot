package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/handler/middleware"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/sponsorblock-plugin"
	"github.com/google/go-github/v52/github"
	"github.com/lavalink-devs/lavalink-bot/commands"
	"github.com/lavalink-devs/lavalink-bot/handlers"
	"github.com/lavalink-devs/lavalink-bot/internal/maven"
	"github.com/lavalink-devs/lavalink-bot/lavalinkbot"
	"github.com/mattn/go-colorable"
	"github.com/topi314/tint"
)

func main() {
	path := flag.String("config", "config.yml", "path to config.yml")
	flag.Parse()

	cfg, err := lavalinkbot.ReadConfig(*path)
	if err != nil {
		slog.Error("failed to read config", tint.Err(err))
		os.Exit(-1)
	}
	setupLogger(cfg.Log)
	slog.Info("starting lavalink-bot...", slog.String("disgo_version", disgo.Version), slog.String("disgolink_version", disgolink.Version))
	slog.Info("Config", slog.String("path", *path), slog.String("config", cfg.String()))

	b := &lavalinkbot.Bot{
		Cfg:    cfg,
		GitHub: github.NewClient(nil),
		Maven: maven.New(&http.Client{
			Timeout: 10 * time.Second,
		}),
		MusicQueue: lavalinkbot.NewPlayerManager(),
	}

	cmds := &commands.Commands{Bot: b}
	r := handler.New()
	r.Use(middleware.Go)
	r.Command("/info/bot", cmds.InfoBot)
	r.Command("/info/lavalink", cmds.InfoLavalink)
	r.Command("/latest", cmds.Latest)
	r.Autocomplete("/latest", cmds.LatestAutocomplete)
	r.Command("/decode", cmds.Decode)

	r.Route("/music", func(r handler.Router) {
		r.Command("/play", cmds.Play)
		r.Command("/tts", cmds.TTS)
		r.Command("/lyrics", cmds.Lyrics)
		r.Autocomplete("/play", cmds.PlayAutocomplete)
		r.Group(func(r handler.Router) {
			r.Use(cmds.RequirePlayer)

			r.Command("/stop", cmds.Stop)
			r.Command("/disconnect", cmds.Disconnect)
			r.Command("/skip", cmds.Skip)
			r.Command("/pause", cmds.Pause)
			r.Command("/resume", cmds.Resume)
			r.Command("/seek", cmds.Seek)
			r.Command("/volume", cmds.Volume)
			r.Command("/shuffle", cmds.Shuffle)
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
			r.Command("/sponsorblock/show", cmds.ShowSponsorblock)
			r.Command("/sponsorblock/set", cmds.SetSponsorblock)
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
		slog.Error("failed to create disgo client", tint.Err(err))
		os.Exit(-1)
	}

	if err = handler.SyncCommands(b.Client, commands.CommandCreates, b.Cfg.Bot.GuildIDs); err != nil {
		slog.Error("failed to sync commands", tint.Err(err))
	}

	sponsorblockPlugin := sponsorblock.New()
	if b.Lavalink = disgolink.New(b.Client.ApplicationID(),
		disgolink.WithPlugins(sponsorblockPlugin),
		disgolink.WithListenerFunc(hdlr.OnTrackStart),
		disgolink.WithListenerFunc(hdlr.OnTrackEnd),
		disgolink.WithListenerFunc(hdlr.OnTrackException),
		disgolink.WithListenerFunc(hdlr.OnTrackStuck),
		disgolink.WithListenerFunc(hdlr.OnWebSocketClosed),
		disgolink.WithListenerFunc(hdlr.OnUnknownEvent),
		disgolink.WithListenerFunc(hdlr.OnSegmentsLoaded),
		disgolink.WithListenerFunc(hdlr.OndSegmentSkipped),
		disgolink.WithListenerFunc(hdlr.OnChaptersLoaded),
		disgolink.WithListenerFunc(hdlr.OnChapterStarted),
	); err != nil {
		slog.Error("failed to create disgolink client", tint.Err(err))
		os.Exit(-1)
	}

	if err = b.Start(); err != nil {
		slog.Error("failed to start bot", tint.Err(err))
		os.Exit(-1)
	}
	defer b.Stop()

	slog.Info("lavalink-bot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

const (
	ansiFaint         = "\033[2m"
	ansiWhiteBold     = "\033[37;1m"
	ansiYellowBold    = "\033[33;1m"
	ansiCyanBold      = "\033[36;1m"
	ansiCyanBoldFaint = "\033[36;1;2m"
	ansiRedFaint      = "\033[31;2m"
	ansiRedBold       = "\033[31;1m"

	ansiRed     = "\033[31m"
	ansiYellow  = "\033[33m"
	ansiGreen   = "\033[32m"
	ansiMagenta = "\033[35m"
)

func setupLogger(cfg lavalinkbot.LogConfig) {
	var sHandler slog.Handler
	switch cfg.Format {
	case "json":
		sHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: cfg.AddSource,
			Level:     cfg.Level,
		})

	case "text":
		sHandler = tint.NewHandler(colorable.NewColorable(os.Stdout), &tint.Options{
			AddSource: cfg.AddSource,
			Level:     cfg.Level,
			NoColor:   cfg.NoColor,
			LevelColors: map[slog.Level]string{
				slog.LevelDebug: ansiMagenta,
				slog.LevelInfo:  ansiGreen,
				slog.LevelWarn:  ansiYellow,
				slog.LevelError: ansiRed,
			},
			Colors: map[tint.Kind]string{
				tint.KindTime:            ansiYellowBold,
				tint.KindSourceFile:      ansiCyanBold,
				tint.KindSourceSeparator: ansiCyanBoldFaint,
				tint.KindSourceLine:      ansiCyanBold,
				tint.KindMessage:         ansiWhiteBold,
				tint.KindKey:             ansiFaint,
				tint.KindSeparator:       ansiFaint,
				tint.KindValue:           ansiWhiteBold,
				tint.KindErrorKey:        ansiRedFaint,
				tint.KindErrorSeparator:  ansiFaint,
				tint.KindErrorValue:      ansiRedBold,
			},
		})
	default:
		slog.Error("Unknown log format", slog.String("format", cfg.Format))
		os.Exit(-1)
	}
	slog.SetDefault(slog.New(sHandler))
}

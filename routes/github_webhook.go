package routes

import (
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"github.com/google/go-github/v52/github"
	"github.com/topi314/tint"

	"github.com/lavalink-devs/lavalink-bot/lavalinkbot"
)

var (
	markdownHeaderRegex            = regexp.MustCompile(`[ \t]*#+[ \t]+([^\r\n]+)`)
	markdownBulletRegex            = regexp.MustCompile(`([ \t]*)[*|-][ \t]+([^\r\n]+)`)
	markdownCheckBoxCheckedRegex   = regexp.MustCompile(`([ \t]*)[*|-][ \t]{0,4}\[x][ \t]+([^\r\n]+)`)
	markdownCheckBoxUncheckedRegex = regexp.MustCompile(`([ \t]*)[*|-][ \t]{0,4}\[ ][ \t]+([^\r\n]+)`)
	prURLRegex                     = regexp.MustCompile(`https?://github\.com/(\w+/\w+)/pull/(\d+)`)
	prNumberRegex                  = regexp.MustCompile(`#(\d+)`)
	commitURLRegex                 = regexp.MustCompile(`https?://github\.com/\w+/\w+/commit/([a-f\d]{7})[a-f\d]+`)
	mentionRegex                   = regexp.MustCompile(`@([a-zA-Z0-9-]+)`)
)

func HandleGithubWebhook(b *lavalinkbot.Bot) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := github.ValidatePayload(r, []byte(b.Cfg.GitHub.WebhookSecret))
		if err != nil {
			slog.Error("Failed to validate payload", tint.Err(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		event, err := github.ParseWebHook(github.WebHookType(r), payload)
		if err != nil {
			slog.Error("Failed to parse webhook", tint.Err(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		switch e := event.(type) {
		case *github.ReleaseEvent:
			err = processReleaseEvent(b, e)
		}
		if err != nil {
			slog.Error("Failed to process webhook", tint.Err(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func processReleaseEvent(b *lavalinkbot.Bot, e *github.ReleaseEvent) error {
	if e.GetAction() != "published" {
		return nil
	}

	repo := e.GetRepo().GetName()
	fullName := e.GetRepo().GetFullName()

	cfg, ok := b.Cfg.GitHub.Releases[fullName]
	if !ok {
		return fmt.Errorf("no config found for %s", fullName)
	}

	webhookClient, ok := b.Webhooks[fullName]
	if !ok {
		webhookClient = webhook.New(cfg.WebhookID, cfg.WebhookToken)
		b.Webhooks[fullName] = webhookClient
	}

	message := parseMarkdown(e.GetRelease().GetBody(), fullName)
	if len(message) > 1024 {
		message = substr(message, 0, 1024)
		if index := strings.LastIndex(message, "\n"); index != -1 {
			message = message[:index]
		}
		message += "\n…"
	}

	msg, err := webhookClient.CreateMessage(discord.NewWebhookMessageCreateBuilder().
		SetContent(discord.RoleMention(cfg.PingRole)).
		SetEmbeds(discord.NewEmbedBuilder().
			SetAuthor(
				fmt.Sprintf("%s version %s has been released", repo, e.Release.GetTagName()),
				e.GetRelease().GetHTMLURL(),
				e.GetRepo().GetOwner().GetAvatarURL(),
			).
			SetDescription(message).
			SetColor(0xff624a).
			SetFooter("Release by "+e.GetRelease().GetAuthor().GetLogin(), e.GetRelease().GetAuthor().GetAvatarURL()).
			SetTimestamp(e.GetRelease().GetCreatedAt().Time).
			Build(),
		).
		SetAvatarURL(e.GetRepo().GetOwner().GetAvatarURL()).
		Build(),
	)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	_, err = b.Client.Rest().CrosspostMessage(msg.ChannelID, msg.ID)
	if err != nil {
		return fmt.Errorf("failed to crosspost message: %w", err)
	}
	return nil
}

func substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}

func parseMarkdown(text string, repo string) string {
	text = markdownCheckBoxCheckedRegex.ReplaceAllString(text, "$1:ballot_box_with_check: $2")
	text = markdownCheckBoxUncheckedRegex.ReplaceAllString(text, "$1:white_square_button: $2")
	text = markdownHeaderRegex.ReplaceAllString(text, "**$1**")
	text = markdownBulletRegex.ReplaceAllString(text, "$1• $2")
	text = prURLRegex.ReplaceAllString(text, "[$1#$2]($0)")
	text = prNumberRegex.ReplaceAllString(text, "[#$1](https://github.com/"+repo+"/pull/$1)")
	text = commitURLRegex.ReplaceAllString(text, "[`$1`]($0)")
	text = mentionRegex.ReplaceAllString(text, "[@$1](https://github.com/$1)")
	return text
}

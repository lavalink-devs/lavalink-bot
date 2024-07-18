package commands

import (
	"github.com/disgoorg/disgo/discord"

	"github.com/lavalink-devs/lavalink-bot/lavalinkbot"
)

var CommandCreates = []discord.ApplicationCommandCreate{
	decode,
	resolve,
	info,
	latest,
	music,
	read,
}

type Commands struct {
	*lavalinkbot.Bot
}

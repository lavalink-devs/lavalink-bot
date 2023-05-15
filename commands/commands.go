package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/lavalink-devs/lavalink-bot/lavalinkbot"
)

var Commands = []discord.ApplicationCommandCreate{
	decode,
	info,
	latest,
	music,
}

type Cmds struct {
	*lavalinkbot.Bot
}

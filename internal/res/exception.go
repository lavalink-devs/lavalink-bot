package res

import (
	"bytes"
	"fmt"

	"github.com/disgoorg/disgo/discord"
)

func NewExceptionFile(causeStackTrace string) *discord.File {
	return &discord.File{
		Name:   "exception.txt",
		Reader: bytes.NewReader([]byte(fmt.Sprintf("```%s```", causeStackTrace))),
	}
}

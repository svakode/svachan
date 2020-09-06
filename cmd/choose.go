package cmd

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"strings"

	"github.com/svakode/svachan/dictionary"
	"github.com/svakode/svachan/framework"
	"github.com/svakode/svachan/helper"
)

func ChooseCommand(ctx *framework.Context) *discordgo.Message {
	args := strings.Join(ctx.Discord.GetArgs(), " ")
	choices := strings.Split(args, ",")

	if len(choices) == 1 {
		return ctx.Discord.Reply(dictionary.NoChoicesMessage)
	}

	result := rand.Intn(len(choices))
	choice := strings.Trim(choices[result], " ")
	message := helper.GetRandomMessage(dictionary.ChooseMessages)

	return ctx.Discord.Reply(fmt.Sprintf(message, choice))
}

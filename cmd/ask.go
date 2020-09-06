package cmd

import (
	"github.com/bwmarrin/discordgo"
	"github.com/svakode/svachan/dictionary"
	"github.com/svakode/svachan/framework"
	"github.com/svakode/svachan/helper"
)

func AskCommand(ctx *framework.Context) *discordgo.Message {
	return ctx.Discord.Reply(helper.GetRandomMessage(dictionary.AskMessages))
}

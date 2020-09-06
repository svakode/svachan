package cmd

import (
	"github.com/bwmarrin/discordgo"
	"github.com/svakode/svachan/dictionary"
	"github.com/svakode/svachan/framework"
)

func HelpCommand(ctx *framework.Context) *discordgo.Message {
	return ctx.Discord.Help(dictionary.HelpMessage)
}

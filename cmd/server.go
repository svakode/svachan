package cmd

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/svakode/svachan/dictionary"
	"github.com/svakode/svachan/framework"
	"github.com/matishsiao/goInfo"
)

func ServerCommand(ctx *framework.Context) *discordgo.Message {
	osInfo := goInfo.GetInfo()
	memoryInfo, err := ctx.Server.Memory.Get()
	if err != nil {
		return ctx.Discord.Reply(dictionary.GeneralError)
	}

	cpuInfo, err := ctx.Server.CPU.Get()
	if err != nil {
		return ctx.Discord.Reply(dictionary.GeneralError)
	}

	diskInfo, err := ctx.Server.Disk.Get()
	if err != nil {
		return ctx.Discord.Reply(dictionary.GeneralError)
	}

	content := fmt.Sprintf("**OS**: %s\n", osInfo.OS)
	content += fmt.Sprintf("**Platform**: %s\n", osInfo.Platform)
	content += fmt.Sprintf("**CPUs**: %d\n", osInfo.CPUs)

	msg := ctx.Discord.GetEmbedLayout()
	msg.Title = "Server Status"
	msg.Description = content
	msg.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: "https://gbf.wiki/images/2/22/Vacation_Slip_square.jpg",
	}
	msg.Fields = []*discordgo.MessageEmbedField{
		{
			Name:   "CPU",
			Value:  fmt.Sprintf("%.2f%%", cpuInfo.Used),
			Inline: true,
		},
		{
			Name:   fmt.Sprintf("Memory (%.0f%%)", memoryInfo.Percentage),
			Value:  fmt.Sprintf("%.2f/%.0fGB", memoryInfo.Used, memoryInfo.Total),
			Inline: true,
		},
		{
			Name:   fmt.Sprintf("Disk (%.0f%%)", diskInfo.Percentage),
			Value:  fmt.Sprintf("%.2f/%.0fGB", diskInfo.Total - diskInfo.Free, diskInfo.Total),
			Inline: true,
		},
	}

	return ctx.Discord.Embed(msg)
}

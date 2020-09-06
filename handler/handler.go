package handler

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/svakode/svachan/config"
	"github.com/svakode/svachan/framework"
	"github.com/svakode/svachan/helper"
)

// Handler interfaces
type Handler interface {
	CommandHandler(s *discordgo.Session, m *discordgo.MessageCreate)
	VoiceStateUpdateHandler(s *discordgo.Session, v *discordgo.VoiceStateUpdate)
	ReadyHandler(s *discordgo.Session, r *discordgo.Ready)
}

func (h *handler) CommandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	user := m.Author
	if user.ID == s.State.User.ID || user.Bot {
		return
	}

	content := m.Content
	content, ok := helper.ValidateContent(content, config.Prefix())
	if !ok {
		return
	}

	args := strings.Fields(content)
	name := strings.ToLower(args[0])
	command, found := h.ctx.CmdHandler.Get(name)
	if !found {
		return
	}

	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		fmt.Println("Error getting channel,", err)
		return
	}
	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		fmt.Println("Error getting guild,", err)
		return
	}

	discord := framework.NewDiscord(s, guild, channel, user, m)
	h.ctx.Discord = discord
	h.ctx.Discord.SetArgs(args[1:])
	c := *command
	c(h.ctx)
}

func (h *handler) VoiceStateUpdateHandler(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	vc := s.VoiceConnections[v.GuildID]
	if vc != nil && v.UserID == vc.UserID && v.ChannelID == "" {
		voice := h.ctx.Voices.GetByGuild(v.GuildID)
		if voice != nil {
			voice.Stop()
			h.ctx.Voices.DeleteSession(voice.GetChannel())
		}
	}
}

func (h *handler) ReadyHandler(s *discordgo.Session, r *discordgo.Ready) {
	err := s.UpdateStatus(0, config.Status())
	if err != nil {
		fmt.Println("Svachan cannot update status, try again later")
		return
	}
}

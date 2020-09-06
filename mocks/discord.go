package mocks

import (
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/mock"
)

type DiscordService struct {
	mock.Mock
}

func (d *DiscordService) GetEmbedLayout() *discordgo.MessageEmbed {
	res := d.Called()
	return res.Get(0).(*discordgo.MessageEmbed)
}

func (d *DiscordService) GetMessage() *discordgo.MessageCreate {
	res := d.Called()
	return res.Get(0).(*discordgo.MessageCreate)
}

func (d *DiscordService) GetGuild() *discordgo.Guild {
	res := d.Called()
	return res.Get(0).(*discordgo.Guild)
}

func (d *DiscordService) GetChannel() *discordgo.Channel {
	res := d.Called()
	return res.Get(0).(*discordgo.Channel)
}

func (d *DiscordService) GetUser() *discordgo.User {
	res := d.Called()
	return res.Get(0).(*discordgo.User)
}

func (d *DiscordService) SetArgs(args []string) {
	d.Called(args)
	return
}

func (d *DiscordService) GetArgs() []string {
	res := d.Called()
	return res.Get(0).([]string)
}

func (d *DiscordService) SetSession(session *discordgo.Session) {
	d.Called(session)
	return
}

func (d *DiscordService) GetSession() *discordgo.Session {
	res := d.Called()
	return res.Get(0).(*discordgo.Session)
}

func (d *DiscordService) RetrieveUser(id string) (*discordgo.User, error) {
	res := d.Called(id)
	return res.Get(0).(*discordgo.User), res.Error(1)
}

func (d *DiscordService) Reply(content string) *discordgo.Message {
	res := d.Called(content)
	return res.Get(0).(*discordgo.Message)
}

func (d *DiscordService) Embed(content *discordgo.MessageEmbed) *discordgo.Message {
	res := d.Called(content)
	return res.Get(0).(*discordgo.Message)
}

func (d *DiscordService) Help(content string) *discordgo.Message {
	res := d.Called(content)
	return res.Get(0).(*discordgo.Message)
}

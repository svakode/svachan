package mocks

import (
	"github.com/bwmarrin/discordgo"
	"github.com/svakode/svachan/framework"
	"github.com/stretchr/testify/mock"
)

type VoiceManagerService struct {
	mock.Mock
}

type VoiceService struct {
	mock.Mock
}

func (v *VoiceManagerService) JoinVoiceChannel(s *discordgo.Session, u *discordgo.User, g *discordgo.Guild) (framework.VoiceService, error) {
	res := v.Called(s, u, g)
	return res.Get(0).(framework.VoiceService), res.Error(1)
}

func (v *VoiceManagerService) GetByGuild(guildID string) framework.VoiceService {
	res := v.Called(guildID)
	if res.Get(0) == nil {
		return nil
	}
	return res.Get(0).(framework.VoiceService)
}

func (v *VoiceManagerService) DeleteSession(channelID string) {
	v.Called(channelID)
	return
}

func (v *VoiceService) GetStatus() bool {
	res := v.Called()
	return res.Get(0).(bool)
}

func (v *VoiceService) GetGuild() string {
	res := v.Called()
	return res.Get(0).(string)
}

func (v *VoiceService) GetChannel() string {
	res := v.Called()
	return res.Get(0).(string)
}

func (v *VoiceService) GetCurrent() framework.Song {
	res := v.Called()
	return res.Get(0).(framework.Song)
}

func (v *VoiceService) GetQueue() []framework.Song {
	res := v.Called()
	return res.Get(0).([]framework.Song)
}

func (v *VoiceService) InitPlayer(ctx *framework.Context) {
	v.Called(ctx)
	return
}

func (v *VoiceService) Play(ctx *framework.Context, title, url string) {
	v.Called(ctx, title, url)
	return
}

func (v *VoiceService) Push(title, url string) {
	v.Called(title, url)
	return
}

func (v *VoiceService) Stop() {
	v.Called()
	return
}

func (v *VoiceService) Skip() {
	v.Called()
	return
}
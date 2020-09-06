package framework

import (
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"github.com/svakode/svachan/dictionary"
	"github.com/svakode/svachan/helper"
)

type VoiceManagerService interface {
	GetByGuild(guildID string) VoiceService
	JoinVoiceChannel(s *discordgo.Session, u *discordgo.User, g *discordgo.Guild) (VoiceService, error)
	DeleteSession(channelID string)
}

type VoiceService interface {
	GetStatus() bool
	GetGuild() string
	GetChannel() string
	GetCurrent() Song
	GetQueue() []Song

	InitPlayer(ctx *Context)
	Play(ctx *Context, title, url string)
	Push(title, url string)
	Stop()
	Skip()
}

type (
	Song struct {
		Title string
		URL   string
	}

	Voice struct {
		Connection  *discordgo.VoiceConnection
		StopChannel chan bool
		Queue       []Song
		Current     Song
		ChannelID   string
		GuildID     string
		Status      bool
	}

	VoiceManager struct {
		Sessions map[string]VoiceService
	}
)

func newVoice(guildID, channelID string, connection *discordgo.VoiceConnection) VoiceService {
	return &Voice{
		Connection:  connection,
		StopChannel: make(chan bool),
		ChannelID:   channelID,
		GuildID:     guildID,
	}
}

func (v *Voice) InitPlayer(ctx *Context) {
	for {
		if len(v.Queue) == 0 {
			v.Status = false
			v.Stop()
			ctx.Voices.DeleteSession(v.ChannelID)

			message := helper.GetRandomMessage(dictionary.MusicStopMessages)
			ctx.Discord.Reply(message)
			break
		}

		song := v.Queue[0]
		v.Current = song
		v.Queue = append(v.Queue[:0], v.Queue[1:]...)
		v.Status = true

		message := helper.GetRandomMessage(dictionary.MusicPlayMessages)
		ctx.Discord.Reply(fmt.Sprintf(message, song.Title))
		dgvoice.PlayAudioFile(v.Connection, song.URL, v.StopChannel)
	}
}

func (v *Voice) Play(ctx *Context, title, url string) {
	song := Song{
		Title: title,
		URL:   url,
	}

	v.Queue = append([]Song{song}, v.Queue...)

	if v.Status {
		v.StopChannel <- true
		time.Sleep(1 * time.Second)
	} else {
		go v.InitPlayer(ctx)
	}
}

func (v *Voice) Push(title, url string) {
	song := Song{
		Title: title,
		URL:   url,
	}

	v.Queue = append(v.Queue, song)
}

func (v *Voice) Skip() {
	if v.Status {
		v.StopChannel <- true
	}
}

func (v *Voice) Stop() {
	v.Skip()
	v.Queue = []Song{}
	_ = v.Connection.Disconnect()
}

func (v Voice) GetStatus() bool {
	return v.Status
}

func (v Voice) GetQueue() []Song {
	return v.Queue
}

func (v Voice) GetGuild() string {
	return v.GuildID
}

func (v Voice) GetCurrent() Song {
	return v.Current
}

func (v Voice) GetChannel() string {
	return v.ChannelID
}

func NewVoiceManager() VoiceManagerService {
	return &VoiceManager{make(map[string]VoiceService)}
}

func (v VoiceManager) DeleteSession(channelID string) {
	delete(v.Sessions, channelID)
}

func (v VoiceManager) GetByGuild(guildID string) VoiceService {
	for _, voice := range v.Sessions {
		if voice.GetGuild() == guildID {
			return voice
		}
	}

	return nil
}

func (v *VoiceManager) JoinVoiceChannel(s *discordgo.Session, u *discordgo.User, g *discordgo.Guild) (VoiceService, error) {
	for _, state := range g.VoiceStates {
		if state.UserID == u.ID {
			vc, err := s.ChannelVoiceJoin(g.ID, state.ChannelID, false, true)
			if err != nil {
				return nil, errors.New(dictionary.JoinChannelError)
			}

			if v.Sessions[state.ChannelID] != nil {
				return v.Sessions[state.ChannelID], nil
			}

			voice := newVoice(g.ID, state.ChannelID, vc)
			v.Sessions[state.ChannelID] = voice

			return voice, nil
		}
	}

	return nil, errors.New(dictionary.NotConnectedToVoiceChannel)
}

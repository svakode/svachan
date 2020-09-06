package cmd

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/svakode/svachan/dictionary"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/svakode/svachan/framework"
	"github.com/svakode/svachan/mocks"
)

type MusicTestSuite struct {
	suite.Suite
	ctx *framework.Context

	discord      *mocks.DiscordService
	youtube      *mocks.YoutubeService
	voice        *mocks.VoiceService
	voiceManager *mocks.VoiceManagerService
}

func (suite *MusicTestSuite) SetupTest() {
	suite.discord = new(mocks.DiscordService)
	suite.youtube = new(mocks.YoutubeService)
	suite.voice = new(mocks.VoiceService)
	suite.voiceManager = new(mocks.VoiceManagerService)
	suite.ctx = &framework.Context{
		Discord: suite.discord,
		Youtube: suite.youtube,
		Voices:  suite.voiceManager,
	}
}

func (suite *MusicTestSuite) TestMusicFindShouldSuccess() {
	suite.discord.On("GetArgs").Return([]string{"find", "query"})
	suite.youtube.On("GetVideoID", "query").Return("1", nil)
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.CommandNotFoundError)
	suite.discord.AssertCalled(suite.T(), "Reply", mock.Anything)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *MusicTestSuite) TestMusicFindShortcutShouldSuccess() {
	suite.discord.On("GetArgs").Return([]string{"f", "query"})
	suite.youtube.On("GetVideoID", "query").Return("1", nil)
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.CommandNotFoundError)
	suite.discord.AssertCalled(suite.T(), "Reply", mock.Anything)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *MusicTestSuite) TestMusicFindShouldReturnErrorGivenInvalidParameter() {
	suite.discord.On("GetArgs").Return([]string{"f"})
	suite.discord.On("Reply", dictionary.InvalidParameter).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.CommandNotFoundError)
	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.InvalidParameter)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *MusicTestSuite) TestMusicFindShouldReturnErrorWhenGetVideoIDFailed() {
	suite.discord.On("GetArgs").Return([]string{"f", "query"})
	suite.youtube.On("GetVideoID", "query").Return("", errors.New(dictionary.GoogleError))
	suite.discord.On("Reply", dictionary.GoogleError).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.CommandNotFoundError)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.InvalidParameter)
	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.GoogleError)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *MusicTestSuite) TestMusicPlayShouldSuccess() {
	suite.discord.On("GetArgs").Return([]string{"play", "query"})
	suite.youtube.On("GetVideoDownloadURL", "query").Return("link", "title", nil)
	suite.discord.On("GetSession").Return(&discordgo.Session{})
	suite.discord.On("GetUser").Return(&discordgo.User{})
	suite.discord.On("GetGuild").Return(&discordgo.Guild{})
	suite.voiceManager.On("JoinVoiceChannel", &discordgo.Session{}, &discordgo.User{}, &discordgo.Guild{}).Return(suite.voice, nil)
	suite.voice.On("Play", suite.ctx, "title", "link").Return()
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.voice.AssertCalled(suite.T(), "Play", suite.ctx, "title", "link")
	suite.discord.AssertCalled(suite.T(), "Reply", mock.Anything)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.InvalidParameter)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.CommandNotFoundError)

	suite.discord.AssertExpectations(suite.T())
	suite.voice.AssertExpectations(suite.T())
}

func (suite *MusicTestSuite) TestMusicPlayShortcutShouldSuccess() {
	suite.discord.On("GetArgs").Return([]string{"p", "query"})
	suite.youtube.On("GetVideoDownloadURL", "query").Return("link", "title", nil)
	suite.discord.On("GetSession").Return(&discordgo.Session{})
	suite.discord.On("GetUser").Return(&discordgo.User{})
	suite.discord.On("GetGuild").Return(&discordgo.Guild{})
	suite.voiceManager.On("JoinVoiceChannel", &discordgo.Session{}, &discordgo.User{}, &discordgo.Guild{}).Return(suite.voice, nil)
	suite.voice.On("Play", suite.ctx, "title", "link").Return()
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.voice.AssertCalled(suite.T(), "Play", suite.ctx, "title", "link")
	suite.discord.AssertCalled(suite.T(), "Reply", mock.Anything)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.InvalidParameter)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.CommandNotFoundError)

	suite.discord.AssertExpectations(suite.T())
	suite.voice.AssertExpectations(suite.T())
}

func (suite *MusicTestSuite) TestMusicPlayShouldReturnErrorGivenInvalidParameter() {
	suite.discord.On("GetArgs").Return([]string{"p"})
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.InvalidParameter)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.CommandNotFoundError)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *MusicTestSuite) TestMusicPlayShouldReturnErrorWhenGetVideoDownloadURLFailed() {
	suite.discord.On("GetArgs").Return([]string{"play", "query"})
	suite.youtube.On("GetVideoDownloadURL", "query").Return("", "", errors.New(dictionary.GoogleError))
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.GoogleError )
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.InvalidParameter)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.CommandNotFoundError)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *MusicTestSuite) TestMusicPlayShouldReturnErrorWhenJoinVoiceChannelFailed() {
	suite.discord.On("GetArgs").Return([]string{"p", "query"})
	suite.youtube.On("GetVideoDownloadURL", "query").Return("link", "title", nil)
	suite.discord.On("GetSession").Return(&discordgo.Session{})
	suite.discord.On("GetUser").Return(&discordgo.User{})
	suite.discord.On("GetGuild").Return(&discordgo.Guild{})
	suite.voiceManager.On("JoinVoiceChannel", &discordgo.Session{}, &discordgo.User{}, &discordgo.Guild{}).Return(suite.voice, errors.New(dictionary.NotConnectedToVoiceChannel))
	suite.discord.On("Reply", dictionary.NotConnectedToVoiceChannel).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.NotConnectedToVoiceChannel)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.InvalidParameter)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.CommandNotFoundError)

	suite.discord.AssertExpectations(suite.T())
	suite.voice.AssertExpectations(suite.T())
}

func (suite *MusicTestSuite) TestMusicQueueShouldSuccess() {
	suite.discord.On("GetArgs").Return([]string{"queue", "query"})
	suite.youtube.On("GetVideoDownloadURL", "query").Return("link", "title", nil)
	suite.discord.On("GetSession").Return(&discordgo.Session{})
	suite.discord.On("GetUser").Return(&discordgo.User{})
	suite.discord.On("GetGuild").Return(&discordgo.Guild{})
	suite.voiceManager.On("JoinVoiceChannel", &discordgo.Session{}, &discordgo.User{}, &discordgo.Guild{}).Return(suite.voice, nil)
	suite.voice.On("GetStatus").Return(true)
	suite.voice.On("Push", "title", "link").Return()
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.voice.AssertCalled(suite.T(), "Push", "title", "link")
	suite.discord.AssertCalled(suite.T(), "Reply", mock.Anything)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.InvalidParameter)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.CommandNotFoundError)

	suite.discord.AssertExpectations(suite.T())
	suite.voice.AssertExpectations(suite.T())
}

func (suite *MusicTestSuite) TestMusicQueueShortcutShouldSuccess() {
	suite.discord.On("GetArgs").Return([]string{"q", "query"})
	suite.youtube.On("GetVideoDownloadURL", "query").Return("link", "title", nil)
	suite.discord.On("GetSession").Return(&discordgo.Session{})
	suite.discord.On("GetUser").Return(&discordgo.User{})
	suite.discord.On("GetGuild").Return(&discordgo.Guild{})
	suite.voiceManager.On("JoinVoiceChannel", &discordgo.Session{}, &discordgo.User{}, &discordgo.Guild{}).Return(suite.voice, nil)
	suite.voice.On("GetStatus").Return(true)
	suite.voice.On("Push", "title", "link").Return()
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.voice.AssertCalled(suite.T(), "Push", "title", "link")
	suite.discord.AssertCalled(suite.T(), "Reply", mock.Anything)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.InvalidParameter)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.CommandNotFoundError)

	suite.discord.AssertExpectations(suite.T())
	suite.voice.AssertExpectations(suite.T())
}

func (suite *MusicTestSuite) TestMusicQueueShouldPlayWhenNotPlaying() {
	suite.discord.On("GetArgs").Return([]string{"q", "query"})
	suite.youtube.On("GetVideoDownloadURL", "query").Return("link", "title", nil)
	suite.discord.On("GetSession").Return(&discordgo.Session{})
	suite.discord.On("GetUser").Return(&discordgo.User{})
	suite.discord.On("GetGuild").Return(&discordgo.Guild{})
	suite.voiceManager.On("JoinVoiceChannel", &discordgo.Session{}, &discordgo.User{}, &discordgo.Guild{}).Return(suite.voice, nil)
	suite.voice.On("GetStatus").Return(false)
	suite.voice.On("Play", suite.ctx, "title", "link").Return()
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.voice.AssertCalled(suite.T(), "Play", suite.ctx, "title", "link")
	suite.discord.AssertCalled(suite.T(), "Reply", mock.Anything)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.InvalidParameter)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.CommandNotFoundError)

	suite.discord.AssertExpectations(suite.T())
	suite.voice.AssertExpectations(suite.T())
}

func (suite *MusicTestSuite) TestMusicQueueShouldReturnErrorGivenInvalidParameter() {
	suite.discord.On("GetArgs").Return([]string{"q"})
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.InvalidParameter)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.CommandNotFoundError)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *MusicTestSuite) TestMusicQueueShouldReturnErrorWhenGetVideoDownloadURLFailed() {
	suite.discord.On("GetArgs").Return([]string{"q", "query"})
	suite.youtube.On("GetVideoDownloadURL", "query").Return("", "", errors.New(dictionary.GoogleError))
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.GoogleError )
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.InvalidParameter)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.CommandNotFoundError)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *MusicTestSuite) TestMusicQueueShouldReturnErrorWhenJoinVoiceChannelFailed() {
	suite.discord.On("GetArgs").Return([]string{"q", "query"})
	suite.youtube.On("GetVideoDownloadURL", "query").Return("link", "title", nil)
	suite.discord.On("GetSession").Return(&discordgo.Session{})
	suite.discord.On("GetUser").Return(&discordgo.User{})
	suite.discord.On("GetGuild").Return(&discordgo.Guild{})
	suite.voiceManager.On("JoinVoiceChannel", &discordgo.Session{}, &discordgo.User{}, &discordgo.Guild{}).Return(suite.voice, errors.New(dictionary.NotConnectedToVoiceChannel))
	suite.discord.On("Reply", dictionary.NotConnectedToVoiceChannel).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.NotConnectedToVoiceChannel)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.InvalidParameter)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.CommandNotFoundError)

	suite.discord.AssertExpectations(suite.T())
	suite.voice.AssertExpectations(suite.T())
}

func (suite *MusicTestSuite) TestMusicSkipShouldSuccess() {
	suite.discord.On("GetArgs").Return([]string{"skip"})
	suite.discord.On("GetGuild").Return(&discordgo.Guild{ID: "1"})
	suite.voiceManager.On("GetByGuild", "1").Return(suite.voice)
	suite.voice.On("Skip").Return()
	suite.voice.On("GetCurrent").Return(framework.Song{})
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.voice.AssertCalled(suite.T(), "Skip")
	suite.discord.AssertCalled(suite.T(), "Reply", mock.Anything)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.CommandNotFoundError)

	suite.discord.AssertExpectations(suite.T())
	suite.voice.AssertExpectations(suite.T())
}

func (suite *MusicTestSuite) TestMusicSkipShortcutShouldSuccess() {
	suite.discord.On("GetArgs").Return([]string{"s"})
	suite.discord.On("GetGuild").Return(&discordgo.Guild{ID: "1"})
	suite.voiceManager.On("GetByGuild", "1").Return(suite.voice)
	suite.voice.On("Skip").Return()
	suite.voice.On("GetCurrent").Return(framework.Song{})
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.voice.AssertCalled(suite.T(), "Skip")
	suite.discord.AssertCalled(suite.T(), "Reply", mock.Anything)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.CommandNotFoundError)

	suite.discord.AssertExpectations(suite.T())
	suite.voice.AssertExpectations(suite.T())
}

func (suite *MusicTestSuite) TestMusicSkipShouldReturnErrorWhenNotPlayingAnything() {
	suite.discord.On("GetArgs").Return([]string{"s"})
	suite.discord.On("GetGuild").Return(&discordgo.Guild{ID: "1"})
	suite.voiceManager.On("GetByGuild", "1").Return(nil)
	suite.discord.On("Reply", dictionary.NotPlaying).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.NotPlaying)

	suite.discord.AssertExpectations(suite.T())
	suite.voice.AssertExpectations(suite.T())
}

func (suite *MusicTestSuite) TestMusicCloseShouldSuccess() {
	suite.discord.On("GetArgs").Return([]string{"close"})
	suite.discord.On("GetGuild").Return(&discordgo.Guild{ID: "1"})
	suite.voiceManager.On("GetByGuild", "1").Return(suite.voice)
	suite.voice.On("Stop").Return()
	suite.voice.On("GetChannel").Return("channel-id")
	suite.voiceManager.On("DeleteSession", "channel-id").Return()
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.voice.AssertCalled(suite.T(), "Stop")
	suite.voiceManager.AssertCalled(suite.T(), "DeleteSession", "channel-id")
	suite.discord.AssertCalled(suite.T(), "Reply", mock.Anything)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.CommandNotFoundError)

	suite.discord.AssertExpectations(suite.T())
	suite.voice.AssertExpectations(suite.T())
	suite.voiceManager.AssertExpectations(suite.T())
}

func (suite *MusicTestSuite) TestMusicCloseShortcutShouldSuccess() {
	suite.discord.On("GetArgs").Return([]string{"c"})
	suite.discord.On("GetGuild").Return(&discordgo.Guild{ID: "1"})
	suite.voiceManager.On("GetByGuild", "1").Return(suite.voice)
	suite.voice.On("Stop").Return()
	suite.voice.On("GetChannel").Return("channel-id")
	suite.voiceManager.On("DeleteSession", "channel-id").Return()
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.voice.AssertCalled(suite.T(), "Stop")
	suite.voiceManager.AssertCalled(suite.T(), "DeleteSession", "channel-id")
	suite.discord.AssertCalled(suite.T(), "Reply", mock.Anything)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.CommandNotFoundError)

	suite.discord.AssertExpectations(suite.T())
	suite.voice.AssertExpectations(suite.T())
	suite.voiceManager.AssertExpectations(suite.T())
}

func (suite *MusicTestSuite) TestMusicStopShouldReturnErrorWhenNotPlayingAnything() {
	suite.discord.On("GetArgs").Return([]string{"c"})
	suite.discord.On("GetGuild").Return(&discordgo.Guild{ID: "1"})
	suite.voiceManager.On("GetByGuild", "1").Return(nil)
	suite.discord.On("Reply", dictionary.NotPlaying).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.NotPlaying)

	suite.discord.AssertExpectations(suite.T())
	suite.voice.AssertExpectations(suite.T())
}

func (suite *MusicTestSuite) TestMusicListShouldSuccess() {
	suite.discord.On("GetArgs").Return([]string{"list"})
	suite.discord.On("GetEmbedLayout").Return(&discordgo.MessageEmbed{})
	suite.discord.On("GetGuild").Return(&discordgo.Guild{ID: "1"})
	suite.voiceManager.On("GetByGuild", "1").Return(suite.voice)
	suite.voice.On("GetQueue").Return([]framework.Song{
		{
			Title: "title-1",
			URL: "link-1",
		},
		{
			Title: "title-2",
			URL: "link-2",
		},
	})
	suite.discord.On("Embed", mock.Anything).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.voice.AssertCalled(suite.T(), "GetQueue")
	suite.discord.AssertCalled(suite.T(), "Embed", mock.Anything)

	suite.discord.AssertExpectations(suite.T())
	suite.voice.AssertExpectations(suite.T())
}

func (suite *MusicTestSuite) TestMusicListShortcutShouldSuccess() {
	suite.discord.On("GetArgs").Return([]string{"l"})
	suite.discord.On("GetEmbedLayout").Return(&discordgo.MessageEmbed{})
	suite.discord.On("GetGuild").Return(&discordgo.Guild{ID: "1"})
	suite.voiceManager.On("GetByGuild", "1").Return(suite.voice)
	suite.voice.On("GetQueue").Return([]framework.Song{
		{
			Title: "title-1",
			URL: "link-1",
		},
		{
			Title: "title-2",
			URL: "link-2",
		},
	})
	suite.discord.On("Embed", mock.Anything).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.voice.AssertCalled(suite.T(), "GetQueue")
	suite.discord.AssertCalled(suite.T(), "Embed", mock.Anything)

	suite.discord.AssertExpectations(suite.T())
	suite.voice.AssertExpectations(suite.T())
}

func (suite *MusicTestSuite) TestMusicListShouldReturnErrorWhenNotPlaying() {
	suite.discord.On("GetArgs").Return([]string{"l"})
	suite.discord.On("GetEmbedLayout").Return(&discordgo.MessageEmbed{})
	suite.discord.On("GetGuild").Return(&discordgo.Guild{ID: "1"})
	suite.voiceManager.On("GetByGuild", "1").Return(nil)
	suite.discord.On("Reply", dictionary.NotPlaying).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.NotPlaying)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *MusicTestSuite) TestMusicListShouldReturnErrorWhenQueueIsEmpty() {
	suite.discord.On("GetArgs").Return([]string{"l"})
	suite.discord.On("GetEmbedLayout").Return(&discordgo.MessageEmbed{})
	suite.discord.On("GetGuild").Return(&discordgo.Guild{ID: "1"})
	suite.voiceManager.On("GetByGuild", "1").Return(suite.voice)
	suite.voice.On("GetQueue").Return([]framework.Song{})
	suite.discord.On("Reply", dictionary.EmptyPlaylist).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.EmptyPlaylist)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *MusicTestSuite) TestMusicHelpShouldSuccess() {
	suite.discord.On("GetArgs").Return([]string{"help"})
	suite.discord.On("Help", dictionary.MusicHelpMessage).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Help", dictionary.MusicHelpMessage)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *MusicTestSuite) TestMusicShouldReturnCommandNotFoundWhenGivenUnknownCommand() {
	suite.discord.On("GetArgs").Return([]string{"unknown-command"})
	suite.discord.On("Reply", dictionary.CommandNotFoundError).Return(&discordgo.Message{})

	MusicCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.CommandNotFoundError)
	suite.discord.AssertExpectations(suite.T())
}

func TestMusicTestSuite(t *testing.T) {
	suite.Run(t, new(MusicTestSuite))
}

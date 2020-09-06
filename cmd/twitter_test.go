package cmd

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/svakode/svachan/dictionary"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/svakode/svachan/framework"
	"github.com/svakode/svachan/mocks"
)

type TwitterTestSuite struct {
	suite.Suite
	ctx *framework.Context
	streamMap map[string][]string

	discord    *mocks.DiscordService
	twitter    *mocks.TwitterService
	repository *mocks.Repository
}

func (suite *TwitterTestSuite) SetupTest() {
	suite.streamMap = make(map[string][]string)
	suite.streamMap["user-id"] = []string{"id", "id2"}

	suite.discord = new(mocks.DiscordService)
	suite.twitter = new(mocks.TwitterService)
	suite.repository = new(mocks.Repository)
	suite.ctx = &framework.Context{
		Discord: suite.discord,
		Twitter: suite.twitter,
		Repository: suite.repository,
	}
}

func (suite *TwitterTestSuite) TestTwitterStreamShouldSuccess() {
	appendedStreamMap := suite.streamMap
	appendedStreamMap["id"] = []string{"channel-id"}

	suite.discord.On("GetArgs").Return([]string{"stream", "username"})
	suite.twitter.On("RetrieveUser", "username").Return(&twitter.User{IDStr: "id"}, nil)
	suite.discord.On("GetChannel").Return(&discordgo.Channel{ID: "channel-id"})
	suite.twitter.On("GetStreamMap", "id").Return([]string{})
	suite.twitter.On("GetStreamMaps").Return(suite.streamMap)
	suite.twitter.On("StopStream").Return()
	suite.repository.On("AddStream", "id", "username", "channel-id").Return(nil)
	suite.twitter.On("SetStreamMaps", appendedStreamMap).Return()
	suite.discord.On("GetSession").Return(&discordgo.Session{})
	suite.twitter.On("OpenStreamConnection", &discordgo.Session{}).Return(nil)
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	TwitterCommand(suite.ctx)

	suite.twitter.AssertCalled(suite.T(), "OpenStreamConnection", &discordgo.Session{})
	suite.twitter.AssertExpectations(suite.T())
}

func (suite *TwitterTestSuite) TestTwitterStreamShouldReturnErrorGivenInvalidParameter() {
	suite.discord.On("GetArgs").Return([]string{"stream"})
	suite.discord.On("Reply", dictionary.InvalidParameter).Return(&discordgo.Message{})

	TwitterCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.InvalidParameter)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *TwitterTestSuite) TestTwitterStreamShouldReturnErrorWhenRetrieveUser() {
	suite.discord.On("GetArgs").Return([]string{"stream", "username"})
	suite.twitter.On("RetrieveUser", "username").Return(&twitter.User{}, errors.New(dictionary.TwitterError))
	suite.discord.On("Reply", dictionary.TwitterError).Return(&discordgo.Message{})

	TwitterCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.TwitterError)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *TwitterTestSuite) TestTwitterStreamShouldReturnErrorWhenAlreadyStreaming() {
	alreadyStreamingErr := fmt.Sprintf(dictionary.AlreadyStreamingError, "username")
	suite.discord.On("GetArgs").Return([]string{"stream", "username"})
	suite.twitter.On("RetrieveUser", "username").Return(&twitter.User{IDStr: "id"}, nil)
	suite.discord.On("GetChannel").Return(&discordgo.Channel{ID: "id"})
	suite.twitter.On("GetStreamMap", "id").Return([]string{"id"})
	suite.discord.On("Reply", alreadyStreamingErr).Return(&discordgo.Message{})

	TwitterCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", alreadyStreamingErr)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *TwitterTestSuite) TestTwitterStreamShouldReturnErrorWhenAddStreamError() {
	suite.discord.On("GetArgs").Return([]string{"stream", "username"})
	suite.twitter.On("RetrieveUser", "username").Return(&twitter.User{IDStr: "id"}, nil)
	suite.discord.On("GetChannel").Return(&discordgo.Channel{ID: "channel-id"})
	suite.twitter.On("GetStreamMap", "id").Return([]string{})
	suite.twitter.On("GetStreamMaps").Return(suite.streamMap)
	suite.twitter.On("StopStream").Return()
	suite.repository.On("AddStream", "id", "username", "channel-id").Return(errors.New(dictionary.DBError))
	suite.discord.On("Reply", dictionary.DBError).Return(&discordgo.Message{})

	TwitterCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.DBError)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *TwitterTestSuite) TestTwitterStreamShouldReturnErrorWhenOpenStreamConnectionFailed() {
	appendedStreamMap := suite.streamMap
	appendedStreamMap["id"] = []string{"channel-id"}

	suite.discord.On("GetArgs").Return([]string{"stream", "username"})
	suite.twitter.On("RetrieveUser", "username").Return(&twitter.User{IDStr: "id"}, nil)
	suite.discord.On("GetChannel").Return(&discordgo.Channel{ID: "channel-id"})
	suite.twitter.On("GetStreamMap", "id").Return([]string{})
	suite.twitter.On("GetStreamMaps").Return(suite.streamMap)
	suite.twitter.On("StopStream").Return()
	suite.repository.On("AddStream", "id", "username", "channel-id").Return(nil)
	suite.twitter.On("SetStreamMaps", appendedStreamMap).Return()
	suite.discord.On("GetSession").Return(&discordgo.Session{})
	suite.twitter.On("OpenStreamConnection", &discordgo.Session{}).Return(errors.New(dictionary.TwitterError))
	suite.discord.On("Reply", dictionary.TwitterError).Return(&discordgo.Message{})

	TwitterCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.TwitterError)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *TwitterTestSuite) TestTwitterStopStreamShouldSuccess() {
	removedStreamMap := make(map[string][]string)
	removedStreamMap["user-id"] = []string{"id2"}

	suite.discord.On("GetArgs").Return([]string{"stop-stream", "username"})
	suite.twitter.On("RetrieveUser", "username").Return(&twitter.User{IDStr: "user-id"}, nil)
	suite.discord.On("GetChannel").Return(&discordgo.Channel{ID: "id"})
	suite.twitter.On("GetStreamMap", "user-id").Return(suite.streamMap["user-id"])
	suite.repository.On("RemoveStream", "user-id", "id").Return(nil)
	suite.twitter.On("GetStreamMaps").Return(suite.streamMap)
	suite.twitter.On("SetStreamMaps", removedStreamMap).Return()
	suite.twitter.On("StopStream").Return()
	suite.discord.On("GetSession").Return(&discordgo.Session{})
	suite.twitter.On("OpenStreamConnection", &discordgo.Session{}).Return(nil)
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	TwitterCommand(suite.ctx)

	suite.twitter.AssertCalled(suite.T(), "OpenStreamConnection", &discordgo.Session{})
	suite.twitter.AssertExpectations(suite.T())
}

func (suite *TwitterTestSuite) TestTwitterStopStreamShouldReturnErrorGivenInvalidParameter() {
	suite.discord.On("GetArgs").Return([]string{"stop-stream"})
	suite.discord.On("Reply", dictionary.InvalidParameter).Return(&discordgo.Message{})

	TwitterCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.InvalidParameter)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *TwitterTestSuite) TestTwitterStopStreamShouldReturnErrorIfRetrieveUserFailed() {
	suite.discord.On("GetArgs").Return([]string{"stop-stream", "username"})
	suite.twitter.On("RetrieveUser", "username").Return(&twitter.User{}, errors.New(dictionary.TwitterError))
	suite.discord.On("Reply", dictionary.TwitterError).Return(&discordgo.Message{})

	TwitterCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.TwitterError)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *TwitterTestSuite) TestTwitterStopStreamShouldReturnErrorIfStreamingNotFound() {
	expectedErr := fmt.Sprintf(dictionary.StreamingNotFoundError, "username")

	suite.discord.On("GetArgs").Return([]string{"stop-stream", "username"})
	suite.twitter.On("RetrieveUser", "username").Return(&twitter.User{IDStr: "user-id"}, nil)
	suite.discord.On("GetChannel").Return(&discordgo.Channel{ID: "id"})
	suite.twitter.On("GetStreamMap", "user-id").Return([]string{})
	suite.discord.On("Reply", expectedErr).Return(&discordgo.Message{})

	TwitterCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", expectedErr)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *TwitterTestSuite) TestTwitterStopStreamShouldReturnErrorWhenRemoveStreamFailed() {
	suite.discord.On("GetArgs").Return([]string{"stop-stream", "username"})
	suite.twitter.On("RetrieveUser", "username").Return(&twitter.User{IDStr: "user-id"}, nil)
	suite.discord.On("GetChannel").Return(&discordgo.Channel{ID: "id"})
	suite.twitter.On("GetStreamMap", "user-id").Return(suite.streamMap["user-id"])
	suite.repository.On("RemoveStream", "user-id", "id").Return(errors.New(dictionary.DBError))
	suite.discord.On("Reply", dictionary.DBError).Return(&discordgo.Message{})

	TwitterCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.DBError)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *TwitterTestSuite) TestTwitterStopStreamShouldReturnErrorWhenOpenStreamConnectionFailed() {
	removedStreamMap := make(map[string][]string)
	removedStreamMap["user-id"] = []string{"id2"}

	suite.discord.On("GetArgs").Return([]string{"stop-stream", "username"})
	suite.twitter.On("RetrieveUser", "username").Return(&twitter.User{IDStr: "user-id"}, nil)
	suite.discord.On("GetChannel").Return(&discordgo.Channel{ID: "id"})
	suite.twitter.On("GetStreamMap", "user-id").Return(suite.streamMap["user-id"])
	suite.repository.On("RemoveStream", "user-id", "id").Return(nil)
	suite.twitter.On("GetStreamMaps").Return(suite.streamMap)
	suite.twitter.On("SetStreamMaps", removedStreamMap).Return()
	suite.twitter.On("StopStream").Return()
	suite.discord.On("GetSession").Return(&discordgo.Session{})
	suite.twitter.On("OpenStreamConnection", &discordgo.Session{}).Return(errors.New(dictionary.TwitterError))
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	TwitterCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.TwitterError)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *TwitterTestSuite) TestTwitterListStreamShouldSuccess() {
	suite.discord.On("GetArgs").Return([]string{"list-stream"})
	suite.discord.On("GetChannel").Return(&discordgo.Channel{ID: "channel-id"})
	suite.repository.On("GetStreamsByChannel", "channel-id").Return([]string{"stream-1"}, nil)
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	TwitterCommand(suite.ctx)

	suite.repository.AssertCalled(suite.T(), "GetStreamsByChannel", "channel-id")
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.EmptyListStreamMessages)
	suite.repository.AssertExpectations(suite.T())
	suite.discord.AssertExpectations(suite.T())
}

func (suite *TwitterTestSuite) TestTwitterListStreamShouldReturnErrorWhenGetStreamsByChannelFailed() {
	suite.discord.On("GetArgs").Return([]string{"list-stream"})
	suite.discord.On("GetChannel").Return(&discordgo.Channel{ID: "channel-id"})
	suite.repository.On("GetStreamsByChannel", "channel-id").Return([]string{}, errors.New(dictionary.DBError))
	suite.discord.On("Reply", dictionary.DBError).Return(&discordgo.Message{})

	TwitterCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.DBError)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *TwitterTestSuite) TestTwitterListStreamShouldReturnErrorWhenGetStreamsByChannelEmpty() {
	suite.discord.On("GetArgs").Return([]string{"list-stream"})
	suite.discord.On("GetChannel").Return(&discordgo.Channel{ID: "channel-id"})
	suite.repository.On("GetStreamsByChannel", "channel-id").Return([]string{}, nil)
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	TwitterCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", mock.Anything)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *MemberTestSuite) TestTwitterHelpShouldSuccess() {
	suite.discord.On("GetArgs").Return([]string{"help"})
	suite.discord.On("Help", dictionary.TweetHelpMessage).Return(&discordgo.Message{})

	TwitterCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Help", dictionary.TweetHelpMessage)
	suite.discord.AssertExpectations(suite.T())
}

func TestTwitterTestSuite(t *testing.T) {
	suite.Run(t, new(TwitterTestSuite))
}

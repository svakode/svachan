package cmd

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/svakode/svachan/framework"
	"github.com/svakode/svachan/mocks"
)

type AskTestSuite struct {
	suite.Suite
	ctx *framework.Context

	discord *mocks.DiscordService
}

func (suite *AskTestSuite) SetupTest() {
	suite.discord = new(mocks.DiscordService)
	suite.ctx = &framework.Context{Discord: suite.discord}
}

func (suite *AskTestSuite) TestAskShouldSuccess() {
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	AskCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", mock.Anything)
	suite.discord.AssertExpectations(suite.T())
}

func TestAskTestSuite(t *testing.T) {
	suite.Run(t, new(AskTestSuite))
}

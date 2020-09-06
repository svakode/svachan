package cmd

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/suite"

	"github.com/svakode/svachan/dictionary"
	"github.com/svakode/svachan/framework"
	"github.com/svakode/svachan/mocks"
)

type HelpTestSuite struct {
	suite.Suite
	ctx *framework.Context

	discord *mocks.DiscordService
}

func (suite *HelpTestSuite) SetupTest() {
	suite.discord = new(mocks.DiscordService)
	suite.ctx = &framework.Context{Discord: suite.discord}
}

func (suite *HelpTestSuite) TestHelpShouldSuccess() {
	suite.discord.On("Help", dictionary.HelpMessage).Return(&discordgo.Message{})

	HelpCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Help", dictionary.HelpMessage)
	suite.discord.AssertExpectations(suite.T())
}

func TestHelpTestSuite(t *testing.T) {
	suite.Run(t, new(HelpTestSuite))
}

package cmd

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/svakode/svachan/dictionary"
	"github.com/svakode/svachan/framework"
	"github.com/svakode/svachan/mocks"
)

type ChooseTestSuite struct {
	suite.Suite
	ctx *framework.Context

	discord *mocks.DiscordService
}

func (suite *ChooseTestSuite) SetupTest() {
	suite.discord = new(mocks.DiscordService)
	suite.ctx = &framework.Context{Discord: suite.discord}
}

func (suite *ChooseTestSuite) TestChooseShouldSuccess() {
	suite.discord.On("GetArgs").Return([]string{"choice", "1,", "choice", "2"})
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	ChooseCommand(suite.ctx)

	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.NoChoicesMessage)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *ChooseTestSuite) TestChooseShouldSendNoChoiceMessageWhenGivenOnlyOneChoice() {
	suite.discord.On("GetArgs").Return([]string{"choice"})
	suite.discord.On("Reply", dictionary.NoChoicesMessage).Return(&discordgo.Message{})

	ChooseCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.NoChoicesMessage)
	suite.discord.AssertExpectations(suite.T())
}

func TestChooseTestSuite(t *testing.T) {
	suite.Run(t, new(ChooseTestSuite))
}

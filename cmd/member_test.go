package cmd

import (
	"errors"
	"github.com/svakode/svachan/constant"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/suite"

	"github.com/svakode/svachan/dictionary"
	"github.com/svakode/svachan/framework"
	"github.com/svakode/svachan/mocks"
)

type MemberTestSuite struct {
	suite.Suite
	ctx *framework.Context

	guild *discordgo.Guild

	discord *mocks.DiscordService
	repository *mocks.Repository
}

func (suite *MemberTestSuite) SetupTest() {
	suite.guild = &discordgo.Guild{
		Members: []*discordgo.Member{
			{
				User: &discordgo.User{
					ID:  "1",
					Bot: false,
				},
				Roles: []string{"role"},
			},
			{
				User: &discordgo.User{
					ID:  "2",
					Bot: false,
				},
				Roles: []string{"role-2"},
			},
			{
				User: &discordgo.User{
					ID:  "3",
					Bot: true,
				},
			},
		},
	}

	suite.discord = new(mocks.DiscordService)
	suite.repository = new(mocks.Repository)
	suite.ctx = &framework.Context{
		Discord: suite.discord,
		Repository: suite.repository,
	}
}

func (suite *MemberTestSuite) TestMemberHelpShouldSuccess() {
	suite.discord.On("GetArgs").Return([]string{"help"})
	suite.discord.On("Help", dictionary.MemberHelpMessage).Return(&discordgo.Message{})

	MemberCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Help", dictionary.MemberHelpMessage)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *MemberTestSuite) TestMemberSetEmailShouldSuccess() {
	suite.discord.On("GetArgs").Return([]string{constant.MemberSetEmailCommand, "username", "email"})
	suite.repository.On("SetMemberEmail", "username", "email").Return(nil)
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	MemberCommand(suite.ctx)

	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.DBError)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.InvalidParameter)
	suite.discord.AssertCalled(suite.T(), "Reply", mock.Anything)
	suite.discord.AssertExpectations(suite.T())
	suite.repository.AssertExpectations(suite.T())
}

func (suite *MemberTestSuite) TestMemberSetEmailShouldReturnErrorGivenInvalidParameter() {
	suite.discord.On("GetArgs").Return([]string{constant.MemberSetEmailCommand, "username"})
	suite.discord.On("Reply", dictionary.InvalidParameter).Return(&discordgo.Message{})

	MemberCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.InvalidParameter)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *MemberTestSuite) TestMemberSetEmailShouldReturnErrorWhenSetMemberEmailFailed() {
	suite.discord.On("GetArgs").Return([]string{constant.MemberSetEmailCommand, "username", "email"})
	suite.repository.On("SetMemberEmail", "username", "email").Return(errors.New(dictionary.DBError))
	suite.discord.On("Reply", dictionary.DBError).Return(&discordgo.Message{})

	MemberCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.DBError)
	suite.discord.AssertExpectations(suite.T())
	suite.repository.AssertExpectations(suite.T())
}

func (suite *MemberTestSuite) TestMemberGetEmailShouldSuccess() {
	suite.discord.On("GetArgs").Return([]string{constant.MemberEmailCommand, "username"})
	suite.repository.On("GetMemberEmail", "username").Return("email", nil)
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	MemberCommand(suite.ctx)

	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.DBError)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.InvalidParameter)
	suite.discord.AssertCalled(suite.T(), "Reply", mock.Anything)
	suite.discord.AssertExpectations(suite.T())
	suite.repository.AssertExpectations(suite.T())
}

func (suite *MemberTestSuite) TestMemberGetEmailShouldReturnErrorGivenInvalidParameter() {
	suite.discord.On("GetArgs").Return([]string{constant.MemberEmailCommand})
	suite.discord.On("Reply", dictionary.InvalidParameter).Return(&discordgo.Message{})

	MemberCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.InvalidParameter)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *MemberTestSuite) TestMemberGetEmailShouldReturnErrorWhenGetMemberEmailFailed() {
	suite.discord.On("GetArgs").Return([]string{constant.MemberEmailCommand, "username"})
	suite.repository.On("GetMemberEmail", "username").Return("", errors.New(dictionary.DBError))
	suite.discord.On("Reply", dictionary.DBError).Return(&discordgo.Message{})

	MemberCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.DBError)
	suite.discord.AssertExpectations(suite.T())
	suite.repository.AssertExpectations(suite.T())
}

func (suite *MemberTestSuite) TestMemberRandomWithoutArgsShouldSuccess() {
	suite.discord.On("GetArgs").Return([]string{constant.MemberRandomCommand})
	suite.discord.On("GetGuild").Return(suite.guild)
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	MemberCommand(suite.ctx)

	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.MemberRandomEmptyMessage)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.GeneralError)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.InvalidParameter)
	suite.discord.AssertCalled(suite.T(), "Reply", mock.Anything)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *MemberTestSuite) TestMemberRandomWithArgsShouldSuccess() {
	suite.discord.On("GetArgs").Return([]string{constant.MemberRandomCommand, "role"})
	suite.discord.On("GetGuild").Return(suite.guild)
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	MemberCommand(suite.ctx)

	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.MemberRandomEmptyMessage)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.GeneralError)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.InvalidParameter)
	suite.discord.AssertCalled(suite.T(), "Reply", mock.Anything)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *MemberTestSuite) TestMemberRandomShouldReturnErrorGivenInvalidParameter() {
	suite.discord.On("GetArgs").Return([]string{constant.MemberRandomCommand, "role", "extra-param"})
	suite.discord.On("Reply", dictionary.InvalidParameter).Return(&discordgo.Message{})

	MemberCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.InvalidParameter)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *MemberTestSuite) TestMemberRandomShouldReturnErrorIfMemberEmpty() {
	suite.discord.On("GetArgs").Return([]string{constant.MemberRandomCommand})
	suite.discord.On("GetGuild").Return(&discordgo.Guild{})
	suite.discord.On("Reply", dictionary.GeneralError).Return(&discordgo.Message{})

	MemberCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.GeneralError)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *MemberTestSuite) TestMemberRandomShouldReturnErrorIfMemberRandomEmpty() {
	suite.discord.On("GetArgs").Return([]string{constant.MemberRandomCommand, "unknown-role"})
	suite.discord.On("GetGuild").Return(suite.guild)
	suite.discord.On("Reply", dictionary.MemberRandomEmptyMessage).Return(&discordgo.Message{})

	MemberCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.MemberRandomEmptyMessage)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *MemberTestSuite) TestMemberShouldReturnErrorWhenCommandNotFound() {
	suite.discord.On("GetArgs").Return([]string{"unknown-command"})
	suite.discord.On("Reply", dictionary.CommandNotFoundError).Return(&discordgo.Message{})

	MemberCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.CommandNotFoundError)
	suite.discord.AssertExpectations(suite.T())
}

func TestMemberTestSuite(t *testing.T) {
	suite.Run(t, new(MemberTestSuite))
}

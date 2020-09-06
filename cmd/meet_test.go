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

type MeetTestSuite struct {
	suite.Suite
	ctx *framework.Context

	guild *discordgo.Guild
	user, unknownUser  *discordgo.User

	discord    *mocks.DiscordService
	calendar   *mocks.CalendarService
	repository *mocks.Repository
}

func (suite *MeetTestSuite) SetupTest() {
	suite.user = &discordgo.User{
		ID:  "1",
		Bot: false,
	}
	suite.unknownUser = &discordgo.User{
		ID:  "2",
		Bot: false,
	}
	suite.guild = &discordgo.Guild{
		ID: "123",
		Members: []*discordgo.Member{
			{
				User: suite.user,
			},
			{
				User: suite.unknownUser,
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
	suite.calendar = new(mocks.CalendarService)
	suite.repository = new(mocks.Repository)
	suite.ctx = &framework.Context{
		Discord:    suite.discord,
		Calendar:   suite.calendar,
		Repository: suite.repository,
	}
}

func (suite *MeetTestSuite) TestMeetWithArgsShouldSuccess() {
	params := []string{"@user1"}
	attendees := make([]interface{}, 2)
	attendees[0] = suite.unknownUser.ID
	attendees[1] = suite.user.ID
	email := "user@gmail.com"
	userMap := make(map[string]string)
	userMap["1"] = email

	suite.discord.On("GetArgs").Return(params)
	suite.discord.On("GetGuild").Return(suite.guild)
	suite.discord.On("GetUser").Return(suite.user)
	suite.calendar.On("RetrieveUsernames", params, suite.guild.Members).Return([]string{suite.unknownUser.ID})
	suite.discord.On("RetrieveUser", suite.user.ID).Return(suite.user, nil)
	suite.discord.On("RetrieveUser", suite.unknownUser.ID).Return(suite.unknownUser, nil)
	suite.repository.On("GetMembersEmail", attendees).Return(userMap, nil)
	suite.calendar.On("AppendAttendees", email).Return()
	suite.calendar.On("StartMeeting", mock.Anything, mock.Anything).Return("https://meet.google.com", nil)
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	MeetCommand(suite.ctx)

	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.GeneralError)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.GoogleError)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.DBError)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.MemberNotFoundMessage)
	suite.discord.AssertCalled(suite.T(), "Reply", mock.Anything)
	suite.discord.AssertExpectations(suite.T())
	suite.calendar.AssertExpectations(suite.T())
	suite.repository.AssertExpectations(suite.T())
}

func (suite *MeetTestSuite) TestMeetWithoutArgsShouldSuccess() {
	attendees := make([]interface{}, 2)
	attendees[0] = suite.user.ID
	attendees[1] = suite.unknownUser.ID
	email := "user@gmail.com"
	userMap := make(map[string]string)
	userMap["1"] = email

	suite.discord.On("GetArgs").Return([]string{})
	suite.discord.On("GetGuild").Return(suite.guild)
	suite.discord.On("GetUser").Return(suite.user)
	suite.discord.On("RetrieveUser", suite.user.ID).Return(suite.user, nil)
	suite.discord.On("RetrieveUser", suite.unknownUser.ID).Return(suite.unknownUser, nil)
	suite.repository.On("GetMembersEmail", attendees).Return(userMap, nil)
	suite.calendar.On("AppendAttendees", email).Return()
	suite.calendar.On("StartMeeting", mock.Anything, mock.Anything).Return("https://meet.google.com", nil)
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	MeetCommand(suite.ctx)

	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.GeneralError)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.GoogleError)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.DBError)
	suite.discord.AssertNotCalled(suite.T(), "Reply", dictionary.MemberNotFoundMessage)
	suite.discord.AssertCalled(suite.T(), "Reply", mock.Anything)
	suite.discord.AssertExpectations(suite.T())
	suite.calendar.AssertExpectations(suite.T())
	suite.repository.AssertExpectations(suite.T())
}

func (suite *MeetTestSuite) TestMeetShouldReturnErrorWhenFailedRetrievingUser() {
	suite.discord.On("GetArgs").Return([]string{})
	suite.discord.On("GetGuild").Return(suite.guild)
	suite.discord.On("GetUser").Return(suite.user)
	suite.discord.On("RetrieveUser", suite.user.ID).Return(&discordgo.User{}, errors.New("some-error"))
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	MeetCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.GeneralError)
	suite.discord.AssertExpectations(suite.T())
}

func (suite *MeetTestSuite) TestMeetShouldReturnErrorWhenGetMembersEmailFailed() {
	attendees := make([]interface{}, 2)
	attendees[0] = suite.user.ID
	attendees[1] = suite.unknownUser.ID

	suite.discord.On("GetArgs").Return([]string{})
	suite.discord.On("GetGuild").Return(suite.guild)
	suite.discord.On("GetUser").Return(suite.user)
	suite.discord.On("RetrieveUser", suite.user.ID).Return(suite.user, nil)
	suite.discord.On("RetrieveUser", suite.unknownUser.ID).Return(suite.unknownUser, nil)
	suite.repository.On("GetMembersEmail", attendees).Return(map[string]string{}, errors.New(dictionary.DBError))
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	MeetCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.DBError)
	suite.repository.AssertExpectations(suite.T())
	suite.discord.AssertExpectations(suite.T())
}

func (suite *MeetTestSuite) TestMeetShouldReturnErrorWhenFailedOnStartMeeting() {
	attendees := make([]interface{}, 2)
	attendees[0] = suite.user.ID
	attendees[1] = suite.unknownUser.ID
	email := "user@gmail.com"
	userMap := make(map[string]string)
	userMap["1"] = email

	suite.discord.On("GetArgs").Return([]string{})
	suite.discord.On("GetGuild").Return(suite.guild)
	suite.discord.On("GetUser").Return(suite.user)
	suite.discord.On("RetrieveUser", suite.user.ID).Return(suite.user, nil)
	suite.discord.On("RetrieveUser", suite.unknownUser.ID).Return(suite.unknownUser, nil)
	suite.repository.On("GetMembersEmail", attendees).Return(userMap, nil)
	suite.calendar.On("AppendAttendees", email).Return()
	suite.calendar.On("StartMeeting", mock.Anything, mock.Anything).Return("", errors.New("some-error"))
	suite.discord.On("Reply", mock.Anything).Return(&discordgo.Message{})

	MeetCommand(suite.ctx)

	suite.discord.AssertCalled(suite.T(), "Reply", dictionary.GoogleError)
	suite.discord.AssertExpectations(suite.T())
	suite.calendar.AssertExpectations(suite.T())
	suite.repository.AssertExpectations(suite.T())
}

func TestMeetTestSuite(t *testing.T) {
	suite.Run(t, new(MeetTestSuite))
}

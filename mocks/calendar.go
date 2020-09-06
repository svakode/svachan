package mocks

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/mock"
)

type CalendarService struct {
	mock.Mock
}

func (c *CalendarService) RetrieveUsernames(attendees []string, members []*discordgo.Member) []string {
	res := c.Called(attendees, members)
	return res.Get(0).([]string)
}

func (c *CalendarService) AppendAttendees(email string) {
	c.Called(email)
	return
}

func (c *CalendarService) StartMeeting(start, end time.Time) (link string, err error) {
	res := c.Called(start, end)
	return res.Get(0).(string), res.Error(1)
}
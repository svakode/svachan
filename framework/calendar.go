package framework

import (
	"context"
	"io/ioutil"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"

	"github.com/svakode/svachan/config"
	"github.com/svakode/svachan/helper"
	"github.com/svakode/svachan/utils"
)

type CalendarService interface {
	AppendAttendees(email string)
	RetrieveUsernames(attendees []string, members []*discordgo.Member) []string
	StartMeeting(start, end time.Time) (link string, err error)
}

type Calendar struct {
	Client    *calendar.Service
	Attendees []*calendar.EventAttendee
}

func NewCalendar(googleConfig *config.GoogleConfig) CalendarService {
	client := NewCalendarClient(googleConfig.AuthFile, googleConfig.DevEmail)

	return &Calendar{
		Client: client,
	}
}

func NewCalendarClient(filePath, devEmail string) *calendar.Service {
	ctx := context.Background()

	jsonCredentials, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic("cannot read json file")
	}

	calendarConfig, err := google.JWTConfigFromJSON(jsonCredentials, admin.AdminDirectoryResourceCalendarScope, calendar.CalendarScope, calendar.CalendarEventsScope)
	if err != nil {
		panic(err.Error())
	}
	calendarConfig.Subject = devEmail

	ts := calendarConfig.TokenSource(ctx)

	svc, err := calendar.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		panic("cannot establish connection with google calendar")
	}

	return svc
}

func (c *Calendar) AppendAttendees(email string) {
	c.Attendees = append(c.Attendees, &calendar.EventAttendee{Email: email})
}

func (c *Calendar) StartMeeting(start, end time.Time) (link string, err error) {
	event := &calendar.Event{
		Attendees: c.Attendees,
		ConferenceData: &calendar.ConferenceData{
			CreateRequest: &calendar.CreateConferenceRequest{
				RequestId: utils.GetRandomString(10),
			},
		},
		End: &calendar.EventDateTime{
			DateTime: end.Format(time.RFC3339),
			TimeZone: "Asia/Jakarta",
		},
		Start: &calendar.EventDateTime{
			DateTime: start.Format(time.RFC3339),
			TimeZone: "Asia/Jakarta",
		},
		Summary: "Svachan Meeting",
	}

	calendarRequest := c.Client.Events.Insert("primary", event).ConferenceDataVersion(1)
	event, err = calendarRequest.Do()
	if err != nil {
		return
	}

	link = event.HangoutLink
	return
}

func (c *Calendar) RetrieveUsernames(attendees []string, members []*discordgo.Member) (usernames []string) {
	args := strings.Join(attendees, " ")
	params := strings.Split(args, ",")

	for _, param := range params {
		isRole := false
		param = strings.Trim(param, " ")
		param = helper.ReplaceDiscordID(param)
		for _, member := range members {
			roleMatched := false
			if member.User.Bot {
				continue
			}

			if utils.Contains(member.Roles, param) {
				isRole = true
				roleMatched = true
			}

			if roleMatched && !utils.Contains(usernames, member.User.ID) {
				usernames = append(usernames, member.User.ID)
			}
		}

		if !isRole && !utils.Contains(usernames, param) {
			usernames = append(usernames, param)
		}
	}

	return
}

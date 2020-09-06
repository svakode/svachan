package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/svakode/svachan/dictionary"
	"github.com/svakode/svachan/framework"
	"github.com/svakode/svachan/helper"
	"github.com/svakode/svachan/utils"
)

func MeetCommand(ctx *framework.Context) *discordgo.Message {
	var usernames, invited, skipped []string
	userMap := make(map[string]string)
	start := time.Now()
	end := start.Add(1 * time.Hour)
	args := ctx.Discord.GetArgs()
	guild := ctx.Discord.GetGuild()
	user := ctx.Discord.GetUser()

	if len(args) >= 1 {
		usernames = ctx.Calendar.RetrieveUsernames(args, guild.Members)
		if !utils.Contains(usernames, user.ID) {
			usernames = append(usernames, user.ID)
		}
	} else if len(args) == 0 {
		for _, member := range guild.Members {
			if member.User.Bot {
				continue
			}
			usernames = append(usernames, member.User.ID)
		}
	}

	attendees := make([]interface{}, len(usernames))
	for key, username := range usernames {
		attendees[key] = username
		user, err := ctx.Discord.RetrieveUser(username)
		if err != nil {
			return ctx.Discord.Reply(dictionary.GeneralError)
		}
		userMap[username] = user.Mention()
	}

	emailMap, err := ctx.Repository.GetMembersEmail(attendees)
	if err != nil {
		return ctx.Discord.Reply(err.Error())
	}

	for username, email := range emailMap {
		invited = append(invited, userMap[username])
		usernames = utils.RemoveFromSlice(usernames, utils.Find(usernames, username))
		ctx.Calendar.AppendAttendees(email)
	}

	link, err := ctx.Calendar.StartMeeting(start, end)
	if err != nil {
		return ctx.Discord.Reply(dictionary.GoogleError)
	}

	message := fmt.Sprintf(dictionary.MeetDetails, strings.Join(invited, ", "))

	if len(usernames) > 0 {
		for _, username := range usernames {
			skipped = append(skipped, userMap[username])
		}
		message += "\n" + fmt.Sprintf(dictionary.SkippedInvitation, strings.Join(skipped, ", "))
	}

	meetMessage := helper.GetRandomMessage(dictionary.MeetMessages)
	message += "\n\n" + fmt.Sprintf(meetMessage, link)
	return ctx.Discord.Reply(message)
}

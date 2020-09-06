package cmd

import (
	"fmt"
	"math/rand"

	"github.com/bwmarrin/discordgo"

	"github.com/svakode/svachan/constant"
	"github.com/svakode/svachan/dictionary"
	"github.com/svakode/svachan/framework"
	"github.com/svakode/svachan/helper"
	"github.com/svakode/svachan/utils"
)

func MemberCommand(ctx *framework.Context) *discordgo.Message {
	var msg string
	subCmd := ctx.Discord.GetArgs()[0]

	if subCmd == constant.HelpCommand {
		return ctx.Discord.Help(dictionary.MemberHelpMessage)
	}

	if subCmd == constant.MemberSetEmailCommand {
		msg = setMemberEmail(ctx)
	} else if subCmd == constant.MemberEmailCommand {
		msg = getMemberEmail(ctx)
	} else if subCmd == constant.MemberRandomCommand {
		msg = random(ctx)
	} else {
		msg = dictionary.CommandNotFoundError
	}

	return ctx.Discord.Reply(msg)
}

func setMemberEmail(ctx *framework.Context) string {
	args := ctx.Discord.GetArgs()
	if len(args) != 3 {
		return dictionary.InvalidParameter
	}

	username := args[1]
	email := args[2]
	err := ctx.Repository.SetMemberEmail(helper.ReplaceDiscordID(username), email)
	if err != nil {
		return err.Error()
	}

	message := helper.GetRandomMessage(dictionary.MemberSetEmailMessages)
	return fmt.Sprintf(message, email, username)
}

func getMemberEmail(ctx *framework.Context) string {
	args := ctx.Discord.GetArgs()
	if len(args) != 2 {
		return dictionary.InvalidParameter
	}

	username := args[1]
	email, err := ctx.Repository.GetMemberEmail(helper.ReplaceDiscordID(username))
	if err != nil {
		return err.Error()
	}

	message := helper.GetRandomMessage(dictionary.MemberGetEmailMessages)
	return fmt.Sprintf(message, username, email)
}

func random(ctx *framework.Context) string {
	var role string
	args := ctx.Discord.GetArgs()

	if len(args) > 2 {
		return dictionary.InvalidParameter
	} else if len(args) == 2 {
		role = helper.ReplaceDiscordID(args[1])
	}

	users := ctx.Discord.GetGuild().Members
	if len(users) == 0 {
		return dictionary.GeneralError
	}

	var usersFiltered []*discordgo.Member
	for _, user := range users {
		if user.User.Bot {
			continue
		}

		if role != "" && !utils.Contains(user.Roles, role) {
			continue
		}
		usersFiltered = append(usersFiltered, user)
	}
	if len(usersFiltered) == 0 {
		return dictionary.MemberRandomEmptyMessage
	}

	result := rand.Intn(len(usersFiltered))
	message := helper.GetRandomMessage(dictionary.MemberRandomMessages)

	return fmt.Sprintf(message, usersFiltered[result].User.Mention())
}

package cmd

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/svakode/svachan/constant"
	"github.com/svakode/svachan/dictionary"
	"github.com/svakode/svachan/framework"
	"github.com/svakode/svachan/helper"
	"github.com/svakode/svachan/utils"
)

func TwitterCommand(ctx *framework.Context) *discordgo.Message {
	var msg string
	subCmd := ctx.Discord.GetArgs()[0]

	if subCmd == constant.HelpCommand {
		return ctx.Discord.Help(dictionary.TweetHelpMessage)
	}

	if subCmd == constant.TwitterStreamCommand {
		msg = stream(ctx)
	} else if subCmd == constant.TwitterStopStreamCommand {
		msg = stopStream(ctx)
	} else if subCmd == constant.TwitterListStreamCommand {
		msg = listStream(ctx)
	} else {
		msg = dictionary.CommandNotFoundError
	}

	return ctx.Discord.Reply(msg)
}

func stream(ctx *framework.Context) string {
	args := ctx.Discord.GetArgs()
	if len(args) != 2 {
		return dictionary.InvalidParameter
	}

	username := args[1]
	user, err := ctx.Twitter.RetrieveUser(username)
	if err != nil {
		return err.Error()
	}

	channel := ctx.Discord.GetChannel()
	if utils.Find(ctx.Twitter.GetStreamMap(user.IDStr), channel.ID) != -1 {
		return fmt.Sprintf(dictionary.AlreadyStreamingError, username)
	}

	streamMaps := ctx.Twitter.GetStreamMaps()
	if len(helper.GetMapKeys(streamMaps)) > 0 {
		ctx.Twitter.StopStream()
	}

	err = ctx.Repository.AddStream(user.IDStr, username, channel.ID)
	if err != nil {
		return err.Error()
	}

	streamMaps[user.IDStr] = append(streamMaps[user.IDStr], channel.ID)
	ctx.Twitter.SetStreamMaps(streamMaps)

	err = ctx.Twitter.OpenStreamConnection(ctx.Discord.GetSession())
	if err != nil {
		return dictionary.TwitterError
	}

	message := helper.GetRandomMessage(dictionary.StreamMessages)

	return fmt.Sprintf(message, username)
}

func stopStream(ctx *framework.Context) string {
	args := ctx.Discord.GetArgs()
	if len(args) != 2 {
		return dictionary.InvalidParameter
	}

	username := args[1]
	user, err := ctx.Twitter.RetrieveUser(username)
	if err != nil {
		return err.Error()
	}

	channel := ctx.Discord.GetChannel()
	stream := ctx.Twitter.GetStreamMap(user.IDStr)
	idx := utils.Find(stream, channel.ID)
	if idx == -1 {
		return fmt.Sprintf(dictionary.StreamingNotFoundError, username)
	}

	err = ctx.Repository.RemoveStream(user.IDStr, channel.ID)
	if err != nil {
		return err.Error()
	}

	streamMaps := ctx.Twitter.GetStreamMaps()
	streamMaps[user.IDStr] = utils.RemoveFromSlice(ctx.Twitter.GetStreamMap(user.IDStr), idx)
	ctx.Twitter.SetStreamMaps(streamMaps)

	ctx.Twitter.StopStream()
	if len(helper.GetMapKeys(streamMaps)) > 0 {
		err = ctx.Twitter.OpenStreamConnection(ctx.Discord.GetSession())
		if err != nil {
			return fmt.Sprintf(dictionary.TwitterError)
		}
	}

	message := helper.GetRandomMessage(dictionary.StopStreamMessages)

	return fmt.Sprintf(message, username)
}

func listStream(ctx *framework.Context) string {
	users, err := ctx.Repository.GetStreamsByChannel(ctx.Discord.GetChannel().ID)
	if err != nil {
		return err.Error()
	} else if len(users) == 0 {
		return helper.GetRandomMessage(dictionary.EmptyListStreamMessages)
	}

	for idx, user := range users {
		users[idx] = fmt.Sprintf("`%s`", user)
	}

	message := helper.GetRandomMessage(dictionary.ListStreamMessages)
	message = message + strings.Join(users, ", ")

	return message
}

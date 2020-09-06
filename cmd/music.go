package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/svakode/svachan/constant"
	"github.com/svakode/svachan/dictionary"
	"github.com/svakode/svachan/framework"
	"github.com/svakode/svachan/helper"
	"github.com/svakode/svachan/utils"
)

func MusicCommand(ctx *framework.Context) *discordgo.Message {
	var msg string
	subCmd := ctx.Discord.GetArgs()[0]

	if utils.Contains(constant.MusicFindCommand, subCmd) {
		msg = find(ctx)
	} else if utils.Contains(constant.MusicPlayCommand, subCmd) {
		msg = play(ctx)
	} else if utils.Contains(constant.MusicSkipCommand, subCmd) {
		msg = skip(ctx)
	} else if utils.Contains(constant.MusicCloseCommand, subCmd) {
		msg = stop(ctx)
	} else if utils.Contains(constant.MusicQueueCommand, subCmd) {
		msg = queue(ctx)
	} else if utils.Contains(constant.MusicListCommand, subCmd) {
		playlist := ctx.Discord.GetEmbedLayout()
		playlist.Title = "Playlist"
		description, err := list(ctx)
		if err == nil {
			playlist.Description = description
			return ctx.Discord.Embed(playlist)
		}
		msg = err.Error()
	} else if subCmd == constant.HelpCommand {
		return ctx.Discord.Help(dictionary.MusicHelpMessage)
	} else {
		msg = dictionary.CommandNotFoundError
	}

	return ctx.Discord.Reply(msg)
}

func find(ctx *framework.Context) string {
	args := ctx.Discord.GetArgs()
	if len(args) == 1 {
		return dictionary.InvalidParameter
	}

	query := strings.Join(args[1:], " ")
	videoID, err := ctx.Youtube.GetVideoID(query)
	if err != nil {
		return err.Error()
	}

	url := fmt.Sprintf(constant.YoutubeURL, videoID)

	message := helper.GetRandomMessage(dictionary.MusicSearchMessages)
	return fmt.Sprintf(message, url)
}

func play(ctx *framework.Context) string {
	args := ctx.Discord.GetArgs()
	if len(args) == 1 {
		return dictionary.InvalidParameter
	}

	query := strings.Join(ctx.Discord.GetArgs()[1:], " ")
	videoURL, videoTitle, err := ctx.Youtube.GetVideoDownloadURL(query)
	if err != nil {
		return err.Error()
	}

	discord := ctx.Discord.GetSession()
	user := ctx.Discord.GetUser()
	guild := ctx.Discord.GetGuild()
	voice, err := ctx.Voices.JoinVoiceChannel(discord, user, guild)
	if err != nil {
		return err.Error()
	}

	voice.Play(ctx, videoTitle, videoURL)

	return ""
}

func queue(ctx *framework.Context) string {
	args := ctx.Discord.GetArgs()
	if len(args) == 1 {
		return dictionary.InvalidParameter
	}

	query := strings.Join(ctx.Discord.GetArgs()[1:], " ")
	videoURL, videoTitle, err := ctx.Youtube.GetVideoDownloadURL(query)
	if err != nil {
		return err.Error()
	}

	discord := ctx.Discord.GetSession()
	user := ctx.Discord.GetUser()
	guild := ctx.Discord.GetGuild()
	voice, err := ctx.Voices.JoinVoiceChannel(discord, user, guild)
	if err != nil {
		return err.Error()
	}

	var message string
	if voice.GetStatus() {
		voice.Push(videoTitle, videoURL)
		message = helper.GetRandomMessage(dictionary.MusicQueueMessages)
	} else {
		voice.Play(ctx, videoTitle, videoURL)
		return ""
	}

	return fmt.Sprintf(message, videoTitle)
}

func list(ctx *framework.Context) (string, error) {
	voice := ctx.Voices.GetByGuild(ctx.Discord.GetGuild().ID)
	if voice == nil {
		return "", errors.New(dictionary.NotPlaying)
	}

	queue := voice.GetQueue()
	if len(queue) == 0 {
		return "", errors.New(dictionary.EmptyPlaylist)
	}

	playlist := ""
	for idx, song := range queue {
		playlist += fmt.Sprintf("%d. %s\n", idx+1, song.Title)
	}

	return playlist, nil
}

func skip(ctx *framework.Context) string {
	voice := ctx.Voices.GetByGuild(ctx.Discord.GetGuild().ID)
	if voice == nil {
		return dictionary.NotPlaying
	}

	voice.Skip()

	message := helper.GetRandomMessage(dictionary.MusicSkipMessages)
	return fmt.Sprintf(message, voice.GetCurrent().Title)
}

func stop(ctx *framework.Context) string {
	voice := ctx.Voices.GetByGuild(ctx.Discord.GetGuild().ID)
	if voice == nil {
		return dictionary.NotPlaying
	}

	voice.Stop()
	ctx.Voices.DeleteSession(voice.GetChannel())

	return helper.GetRandomMessage(dictionary.MusicStopMessages)
}

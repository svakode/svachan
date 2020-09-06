package framework

import "github.com/svakode/svachan/repository"

type Context struct {
	CmdHandler *CommandHandler

	Discord  DiscordService
	Calendar CalendarService
	Voices   VoiceManagerService
	Youtube  YoutubeService
	Twitter  TwitterService

	Repository repository.Repository
	Server     *Server
}

func NewContext(discord DiscordService, cmdHandler *CommandHandler, repository repository.Repository,
	twitter TwitterService, calendar CalendarService, youtube YoutubeService, voices VoiceManagerService,
	server *Server) *Context {
	ctx := new(Context)
	ctx.Discord = discord
	ctx.CmdHandler = cmdHandler
	ctx.Repository = repository
	ctx.Twitter = twitter
	ctx.Calendar = calendar
	ctx.Youtube = youtube
	ctx.Voices = voices
	ctx.Server = server
	return ctx
}

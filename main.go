package main

import (
	"database/sql"
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/svakode/svachan/cmd"
	"github.com/svakode/svachan/config"
	"github.com/svakode/svachan/constant"
	"github.com/svakode/svachan/framework"
	"github.com/svakode/svachan/handler"
	"github.com/svakode/svachan/repository"
)

var (
	db  *sql.DB
	ctx *framework.Context

	Handler    handler.Handler
	Repository repository.Repository
	Server     *framework.Server
	CmdHandler *framework.CommandHandler

	Discord    framework.DiscordService
	Calendar   framework.CalendarService
	Voices     framework.VoiceManagerService
	Youtube    framework.YoutubeService
	Twitter    framework.TwitterService
)

func init() {
	config.Load()
	db = config.InitDB()

	CmdHandler = framework.NewCommandHandler()
	registerCommands()

	Discord = framework.NewDiscord(&discordgo.Session{}, &discordgo.Guild{}, &discordgo.Channel{}, &discordgo.User{}, &discordgo.MessageCreate{})
	Repository = repository.NewRepository(db)
	Server = framework.NewServer(framework.NewMemory(), framework.NewCPU(), framework.NewDisk())
	Twitter = framework.NewTwitter(config.Twitter())
	Calendar = framework.NewCalendar(config.Google())
	Youtube = framework.NewYoutube(config.Google())
	Voices = framework.NewVoiceManager()

	ctx = framework.NewContext(Discord, CmdHandler, Repository, Twitter, Calendar, Youtube, Voices, Server)
	Handler = handler.NewHandler(ctx)
}

func main() {
	defer db.Close()

	svachan, err := discordgo.New("Bot " + config.DiscordToken())
	if err != nil {
		fmt.Println("Svachan room is locked for now, try again later")
		return
	}

	svachan.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	svachan.AddHandler(Handler.CommandHandler)
	svachan.AddHandler(Handler.VoiceStateUpdateHandler)
	svachan.AddHandler(Handler.ReadyHandler)

	err = svachan.Open()
	if err != nil {
		fmt.Println("Svachan cannot serve you for now")
		return
	}

	ctx.Discord.SetSession(svachan)

	recoverState()

	fmt.Println("Svachan is ready to serve you")
	<-make(chan struct{})
}

func recoverState() {
	if config.Twitter().Ready() {
		err := Twitter.Recover(ctx.Discord.GetSession(), ctx.Repository)
		if err != nil {
			panic("failed to recover twitter state, " + err.Error())
		}
	}
}

func registerCommands() {
	CmdHandler.Register(constant.AskCommand, cmd.AskCommand)
	CmdHandler.Register(constant.ChooseCommand, cmd.ChooseCommand)
	CmdHandler.Register(constant.HelpCommand, cmd.HelpCommand)
	CmdHandler.Register(constant.MeetCommand, cmd.MeetCommand)
	CmdHandler.Register(constant.MemberCommand, cmd.MemberCommand)

	if config.Twitter().Ready() {
		CmdHandler.Register(constant.TwitterCommand, cmd.TwitterCommand)
	}

	if config.Google().Ready() {
		CmdHandler.Register(constant.MusicCommand, cmd.MusicCommand)
		CmdHandler.Register(constant.ServerCommand, cmd.ServerCommand)
	}
}

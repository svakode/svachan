# Svachan

<p align="center"><img src="assets/svachan.png" width="360"></p>

[![Go report](http://goreportcard.com/badge/svakode/svachan)](http://goreportcard.com/report/svakode/svachan) [![Build Status](https://travis-ci.org/svakode/svachan.svg?branch=master)](https://travis-ci.org/svakode/svachan)

Svachan is a [Go](https://golang.org/) Discord Bot. This bot is built on top of [DiscordGo](https://github.com/bwmarrin/discordgo) Multi purposes bot which can be extended 
as much as you want, from handling an event in Discord to adding your own custom command!

## Getting Started

### Installing

This assumes you already have a working Go environment, if not please see
[this page](https://golang.org/doc/install) first.

After you clone your project, you need to run
```
make copy-config
``` 

Open the generated `application.yml` and fill in the required configuration. The configuration is
not required for external parties (e.g Twitter & Google). If empty then we will not register the command to Svachan. 

### Build

You can build the binary file by simply running

```
make build
```

After the build is done, you can run it by 

```
./out/svachan
```

### Handler Extension

You can extend the bot to handle more event by creating the handler function in `handler` package and registering it in the `main.go`. 
Register your handler in the `svachan` variable. For example:

```go
svachan.AddHandler(Handler.CommandHandler)
```

Refer to [this page](https://discord.com/developers/docs/topics/gateway#commands-and-events) for more information about events.

### Command

You can extend the bot to handle more command by creating the handler function in `cmd` package and registering it in the `main.go`. 
Register your handler in the `registerCommand()` function. For example:

```go
CmdHandler.Register(constant.HelpCommand, cmd.HelpCommand)
```

### Built-in Command

| Command     | Trigger     | Description |
| ----------- | ----------- | ----------- |
| Ask         | s.ask       | Answer your yes/no question|
| Choose      | s.choose    | Choose one of given options|
| Help        | s.help      | Showing instruction on how to use |
| Meet        | s.meet      | Initiate a Google Meet with your team |
| Member      | s.member    | Managing your member data in the server |
| Music       | s.music     | Playing music in your server |
| Tweet       | s.tweet     | Listening to users tweet and repost it in the channel |
| Server      | s.server    | Showing server status where Svachan is hosted |

## Credit

Made with :heart: by [Hansen Edrick Harianto](https://github.com/hansenedrickh)
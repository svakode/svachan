# Handler Extension

You can extend the bot to handle more event by creating the handler function in `handler` package and registering it in the `main.go`. 
Register your handler in the `svachan` variable. For example:

```go
svachan.AddHandler(Handler.CommandHandler)
```

Refer to [this page](https://discord.com/developers/docs/topics/gateway#commands-and-events) for more information about events.
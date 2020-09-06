package dictionary

var StreamMessages = []string{
	"started stalking `%s` tweets",
	"listening to `%s` tweets",
}

var StopStreamMessages = []string{
	"stopped stalking `%s` tweets",
	"stop listening to `%s` tweets",
}

var ListStreamMessages = []string{
	"currently stalking ",
}

var EmptyListStreamMessages = []string{
	"currently stalking no one",
}

var StreamActionMessages = []string{
	"saw a new tweet from `%s`: %s",
}

var TweetHelpMessage = "**s.tweet stream <twitter-username>**\n" +
	"Ex: s.tweet stream @Username\n" +
	"Track a specific user's tweet and post it in the channel whenever there is a new tweet\n\n" +
	"**s.tweet stop-stream <twitter-username>**\n" +
	"Ex: s.tweet stop-stream @Username\n" +
	"Stop tracking a specific user's tweet\n\n" +
	"**s.tweet list-stream**\n" +
	"Get a list of tracked user in the channel\n"

var AlreadyStreamingError = "already streaming for `%s`"
var StreamingNotFoundError = "no streaming found for `%s`"

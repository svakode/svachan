package dictionary

var MusicSearchMessages = []string{
	"found this %s",
}

var MusicPlayMessages = []string{
	"playing `%s`",
}

var MusicQueueMessages = []string{
	"queuing `%s`",
}

var MusicSkipMessages = []string{
	"skipped `%s`",
}

var MusicStopMessages = []string{
	"let me know if u want me to play some music again!",
}

var MusicHelpMessage = "**s.music play/p <query>**\n" +
	"Ex: s.music play Westlife\n" +
	"Playing the first song it can found on Youtube for Westlife\n\n" +
	"**s.music queue/q <query>**\n" +
	"Ex: s.music queue Westlife\n" +
	"Queue a song if currently playing another song, otherwise play the song\n\n" +
	"**s.music list/l**\n" +
	"Showing the queue playlist\n\n" +
	"**s.music skip/s**\n" +
	"Skip the current song\n\n" +
	"**s.music close/c**\n" +
	"Stop current player and close the session\n"

var NotConnectedToVoiceChannel = "you are not connected to a voice channel, please connect to one of them before playing any music"
var NotPlaying = "but, i don't play any music now"
var EmptyPlaylist = "it's empty as of now, maybe u want to add some"

package dictionary

var MemberSetEmailMessages = []string{
	"`%s` set to %s",
}

var MemberGetEmailMessages = []string{
	"found it! %s email is `%s`",
}

var MemberRandomMessages = []string{
	"%s is the lucky one",
	"%s is the chosen one",
	"congrats, %s !",
}

var MemberHelpMessage = "**s.member set-email <username> <email>**\n" +
	"Ex: s.member set-email @Username username@gmail.com\n" +
	"Set your email in the bot which can be used for lots of things e.g initiate meeting\n\n" +
	"**s.member email <username>**\n" +
	"Ex: s.member email @Username\n" +
	"Get your saved email in the bot\n\n" +
	"**s.member random <role?>**\n" +
	"Ex: s.member random @role\n" +
	"Get a random member from a role if specified, otherwise will select one random member from the guild\n"

var MemberNotFoundMessage = "i can't find any email related to it, looks like it is not registered"
var MemberRandomEmptyMessage = "cannot find someone who match the criteria"

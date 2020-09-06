package dictionary

// MeetMessages list to be return when executed properly
var MeetMessages = []string{
	"here you go: %s",
	"break a leg: %s",
	"have a great discussion: %s",
}

// MeetDetails is the message to inform who's in the meeting
var MeetDetails = "i have invited %s to the meeting"

// SkippedInvitation is the message to inform who's not in the meeting due to there is no registered email
var SkippedInvitation = "but i can't find %s in my system, please register yourself using **s.member**"

// NoAttendees Message which will be returned if wrong number of parameter given
var NoAttendees = "no attendees specified, please specify attendees(csv) in the first parameter"

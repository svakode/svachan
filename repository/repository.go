package repository

// Repository interfaces
type Repository interface {
	SetMemberEmail(username, email string) error
	GetMemberEmail(username string) (string, error)
	GetMembersEmail(usernames []interface{}) (map[string]string, error)

	AddStream(twitterID, twitterUsername, channelID string) error
	RemoveStream(twitterID, channelID string) error
	GetStreams() (map[string][]string, error)
	GetStreamsByChannel(channelID string) ([]string, error)
}

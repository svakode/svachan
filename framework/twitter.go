package framework

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	"github.com/svakode/svachan/config"
	"github.com/svakode/svachan/dictionary"
	"github.com/svakode/svachan/helper"
	"github.com/svakode/svachan/repository"
)

type TwitterService interface {
	GetStreamMap(userID string) []string
	GetStreamMaps() map[string][]string
	SetStreamMaps(streamMap map[string][]string)
	OpenStreamConnection(s *discordgo.Session) error
	RetrieveUser(username string) (*twitter.User, error)
	Recover(s *discordgo.Session, r repository.Repository) error
	StopStream()
}

type Twitter struct {
	Client        *twitter.Client
	StreamMap     map[string][]string
	StreamChannel *twitter.Stream
}

func NewTwitter(twitterConfig *config.TwitterConfig) TwitterService {
	oauthConfig := oauth1.NewConfig(twitterConfig.ConsumerKey, twitterConfig.ConsumerSecret)
	token := oauth1.NewToken(twitterConfig.AccessKey, twitterConfig.AccessSecret)
	httpClient := oauthConfig.Client(oauth1.NoContext, token)

	return &Twitter{
		Client:    twitter.NewClient(httpClient),
		StreamMap: make(map[string][]string),
	}
}

func (t *Twitter) GetStreamMap(userID string) []string {
	return t.StreamMap[userID]
}

func (t *Twitter) GetStreamMaps() map[string][]string {
	return t.StreamMap
}

func (t *Twitter) SetStreamMaps(streamMap map[string][]string) {
	t.StreamMap = streamMap
}

func (t *Twitter) StopStream() {
	t.StreamChannel.Stop()
}

func (t *Twitter) OpenStreamConnection(s *discordgo.Session) (err error) {
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		if len(t.StreamMap[tweet.User.IDStr]) > 0 {
			tweetURL := fmt.Sprintf("https://twitter.com/%v/status/%v", tweet.User.ScreenName, tweet.ID)

			for _, channel := range t.StreamMap[tweet.User.IDStr] {
				message := helper.GetRandomMessage(dictionary.StreamActionMessages)
				_, err := s.ChannelMessageSend(channel, fmt.Sprintf(message, "@"+tweet.User.ScreenName, tweetURL))
				if err != nil {
					fmt.Println("failed sending message to discord")
				}
			}
		}
	}

	param := &twitter.StreamFilterParams{
		Follow:        helper.GetMapKeys(t.StreamMap),
		StallWarnings: twitter.Bool(true),
	}
	t.StreamChannel, err = t.Client.Streams.Filter(param)
	if err != nil {
		return err
	}

	go demux.HandleChan(t.StreamChannel.Messages)

	return nil
}

func (t *Twitter) RetrieveUser(username string) (*twitter.User, error) {
	screenName := username[1:]
	user, _, err := t.Client.Users.Show(&twitter.UserShowParams{
		ScreenName: screenName,
	})
	if err != nil {
		return nil, fmt.Errorf(dictionary.TwitterError)
	}

	return user, nil
}

func (t *Twitter) Recover(s *discordgo.Session, r repository.Repository) (err error) {
	t.StreamMap, err = r.GetStreams()
	if err != nil {
		return err
	}

	return t.OpenStreamConnection(s)
}

package config

import (
	"github.com/svakode/svachan/utils"
)

type TwitterConfig struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessKey      string
	AccessSecret   string
}

func NewTwitterConfig() *TwitterConfig {
	return &TwitterConfig{
		ConsumerKey:    utils.GetString("TWITTER_API_KEY"),
		ConsumerSecret: utils.GetString("TWITTER_API_SECRET"),
		AccessKey:      utils.GetString("TWITTER_CONSUMER_KEY"),
		AccessSecret:   utils.GetString("TWITTER_CONSUMER_SECRET"),
	}
}

func (t TwitterConfig) Ready() bool {
	return t.ConsumerKey != "" && t.ConsumerSecret != "" && t.AccessKey != "" && t.AccessSecret != ""
}

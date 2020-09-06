package config

import (
	"github.com/svakode/svachan/utils"
)

type GoogleConfig struct {
	AuthFile string
	DevEmail string
	APIKey   string
}

func NewGoogleConfig() *GoogleConfig {
	return &GoogleConfig{
		AuthFile: utils.GetString("GOOGLE_AUTH_FILE"),
		DevEmail: utils.GetString("GOOGLE_DEV_EMAIL"),
		APIKey:   utils.GetString("GOOGLE_API_KEY"),
	}
}

func (g GoogleConfig) Ready() bool {
	return g.AuthFile != "" && g.DevEmail != "" && g.APIKey != ""
}

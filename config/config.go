package config

import (
	"github.com/spf13/viper"

	"github.com/svakode/svachan/utils"
)

type Config struct {
	discordToken string
	prefix       string
	status       string

	twitter      *TwitterConfig
	google       *GoogleConfig
}

var appConfig *Config

func Load() {
	viper.AutomaticEnv()
	viper.SetConfigName("application")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./..")
	viper.SetConfigType("yaml")

	viper.ReadInConfig()

	appConfig = &Config{
		discordToken: utils.FatalGetString("DISCORD_TOKEN"),
		prefix:       utils.FatalGetString("PREFIX"),
		status:       utils.FatalGetString("STATUS"),
		twitter:      NewTwitterConfig(),
		google:       NewGoogleConfig(),
	}
}

func DiscordToken() string {
	return appConfig.discordToken
}

func Prefix() string {
	return appConfig.prefix
}

func Status() string {
	return appConfig.status
}

func Twitter() *TwitterConfig {
	return appConfig.twitter
}

func Google() *GoogleConfig {
	return appConfig.google
}

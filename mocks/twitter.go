package mocks

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/stretchr/testify/mock"

	"github.com/svakode/svachan/repository"
)

type TwitterService struct {
	mock.Mock
}

func (t *TwitterService) GetStreamMap(userID string) []string {
	res := t.Called(userID)
	return res.Get(0).([]string)
}

func (t *TwitterService) GetStreamMaps() map[string][]string {
	res := t.Called()
	return res.Get(0).(map[string][]string)
}

func (t *TwitterService) SetStreamMaps(streamMap map[string][]string) {
	t.Called(streamMap)
	return
}

func (t *TwitterService) OpenStreamConnection(s *discordgo.Session) error {
	res := t.Called(s)
	return res.Error(0)
}

func (t *TwitterService) RetrieveUser(username string) (*twitter.User, error) {
	res := t.Called(username)
	return res.Get(0).(*twitter.User), res.Error(1)
}

func (t *TwitterService) Recover(s *discordgo.Session, r repository.Repository) error {
	res := t.Called(s, r)
	return res.Error(0)
}

func (t *TwitterService) StopStream() {
	t.Called()
	return
}
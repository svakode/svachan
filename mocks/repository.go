package mocks

import (
	"github.com/stretchr/testify/mock"
)

type Repository struct {
	mock.Mock
}

func (r *Repository) SetMemberEmail(username, email string) error {
	res := r.Called(username, email)
	return res.Error(0)
}

func (r *Repository) GetMemberEmail(username string) (string, error) {
	res := r.Called(username)
	return res.Get(0).(string), res.Error(1)
}

func (r *Repository) GetMembersEmail(usernames []interface{}) (map[string]string, error) {
	res := r.Called(usernames)
	return res.Get(0).(map[string]string), res.Error(1)
}

func (r *Repository) AddStream(twitterID, twitterUsername, channelID string) error {
	res := r.Called(twitterID, twitterUsername, channelID)
	return res.Error(0)
}

func (r *Repository) RemoveStream(twitterID, channelID string) error {
	res := r.Called(twitterID, channelID)
	return res.Error(0)
}

func (r *Repository) GetStreams() (map[string][]string, error) {
	res := r.Called()
	return res.Get(0).(map[string][]string), res.Error(1)
}

func (r *Repository) GetStreamsByChannel(channelID string) ([]string, error) {
	res := r.Called(channelID)
	return res.Get(0).([]string), res.Error(1)
}

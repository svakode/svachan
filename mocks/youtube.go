package mocks

import (
	"github.com/stretchr/testify/mock"
)

type YoutubeService struct {
	mock.Mock
}

func (c *YoutubeService) GetVideoID(query string) (string, error) {
	res := c.Called(query)
	return res.Get(0).(string), res.Error(1)
}

func (c *YoutubeService) GetVideoDownloadURL(query string) (string, string, error) {
	res := c.Called(query)
	return res.Get(0).(string), res.Get(1).(string), res.Error(2)
}
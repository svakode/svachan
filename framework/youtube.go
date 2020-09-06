package framework

import (
	"context"
	"errors"

	ytdl "github.com/kkdai/youtube/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"

	"github.com/svakode/svachan/config"
	"github.com/svakode/svachan/dictionary"
)

type YoutubeService interface {
	GetVideoID(query string) (string, error)
	GetVideoDownloadURL(query string) (string, string, error)
}

type Youtube struct {
	Client *youtube.Service
}

func NewYoutube(googleConfig *config.GoogleConfig) YoutubeService {
	client := NewYoutubeClient(googleConfig.APIKey)

	return &Youtube{
		Client: client,
	}
}

func NewYoutubeClient(APIKey string) *youtube.Service {
	ctx := context.Background()

	svc, err := youtube.NewService(ctx, option.WithAPIKey(APIKey))
	if err != nil {
		panic("cannot establish connection with youtube")
	}

	return svc
}

func (y Youtube) GetVideoID(query string) (string, error) {
	call := y.Client.Search.List("id, snippet").
		Q(query).
		MaxResults(1).
		Type("video").
		Order("viewCount")

	videos, err := call.Do()
	if err != nil || len(videos.Items) == 0 {
		return "", errors.New(dictionary.GoogleError)
	}

	return videos.Items[0].Id.VideoId, nil
}

func (y Youtube) GetVideoDownloadURL(query string) (string, string, error) {
	videoID, err := y.GetVideoID(query)
	if err != nil {
		return "", "", err
	}

	yt := ytdl.Client{}
	video, err := yt.GetVideo(videoID)
	if err != nil {
		return "", "", errors.New(dictionary.GoogleError)
	}

	url, err := yt.GetStreamURL(video, &video.Formats[0])
	if err != nil {
		return "", "", errors.New(dictionary.GoogleError)
	}

	return url, video.Title, nil
}

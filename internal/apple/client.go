package apple

import (
	"Test_api/pkg/httpclient"
)

const (
	baseURL   = "https://api.music.apple.com/v1"
	originURL = "https://music.apple.com"
)

func newClient(DevToken, MediaUserToken string) *httpclient.Client {
	return httpclient.New(
		httpclient.WithBaseURL(baseURL),
		httpclient.WithHeader("Authorization", "Bearer "+DevToken),
		httpclient.WithHeader("Music-User-Token", MediaUserToken),
		httpclient.WithHeader("Origin", originURL),
		httpclient.WithHeader("Referer", originURL),
	)
}

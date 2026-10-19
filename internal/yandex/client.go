package yandex

import (
	"Test_api/pkg/httpclient"
)

const (
	baseURL = "https://api.music.yandex.ru"
)

func newClient() *httpclient.Client {
	return httpclient.New(
		httpclient.WithBaseURL(baseURL),
	)
}

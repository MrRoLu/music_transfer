package models

type RespYandexGet struct {
	Result struct {
		Tracks []SongGet `json:"tracks"`
	} `json:"result"`
}

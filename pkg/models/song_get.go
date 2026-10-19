package models

type RespYandexGet struct {
	Result struct {
		Tracks []SongGet `json:"tracks"`
	} `json:"result"`
}

type SongGet struct {
	ID      *uint `json:"id"`
	AlbumID *uint `json:"albumId"`
}

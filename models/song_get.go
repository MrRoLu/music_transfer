package models

type SongGet struct {
	ID      *uint `json:"id"`
	AlbumID *uint `json:"albumId"`
}

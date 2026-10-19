package models

type RespYandexPost struct {
	Result []struct {
		ID      *string  `json:"id"`
		RealID  *string  `json:"realId"`
		Title   *string  `json:"title"`
		Artists []Artist `json:"artists"`
	} `json:"result"`
}

type SongYa struct {
	Title      *string `json:"title"`
	ArtistName *string `json:"artist"`
}

type Artist struct {
	ID   *uint   `json:"id"`
	Name *string `json:"name"`
}

package models

type RespYandexPost struct {
	Result []struct {
		ID      *string  `json:"id"`
		RealId  *string  `json:"realId"`
		Title   *string  `json:"title"`
		Artists []Artist `json:"artists"`
	} `json:"result"`
}

type Artist struct {
	ID   *uint   `json:"id"`
	Name *string `json:"name"`
}

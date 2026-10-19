package models

type StorefrontResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

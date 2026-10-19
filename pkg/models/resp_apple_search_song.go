package models

type CatalogSearchResponse struct {
	Results struct {
		Songs struct {
			Data []struct {
				ID         string `json:"id"`
				Attributes struct {
					Name       string `json:"name"`
					ArtistName string `json:"artistName"`
					AlbumName  string `json:"albumName"`
					ISRC       string `json:"isrc"`
				} `json:"attributes"`
			} `json:"data"`
		} `json:"songs"`
	} `json:"results"`
}

type Song struct {
	ID         string
	Name       string
	ArtistName string
	AlbumName  string
	ISRC       string
}

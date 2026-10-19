package response

import (
	"Test_api/pkg/models"
	"encoding/json"
)

func ParsedCatalogSearchResult(raw []byte) (*[]models.Song, error) {
	var res *models.CatalogSearchResponse

	if err := json.Unmarshal(raw, &res); err != nil {
		return nil, err
	}

	songs := make([]models.Song, len(res.Results.Songs.Data))

	for _, v := range res.Results.Songs.Data {
		songs = append(songs, models.Song{
			ID:         v.ID,
			Name:       v.Attributes.Name,
			ArtistName: v.Attributes.ArtistName,
			AlbumName:  v.Attributes.AlbumName,
			ISRC:       v.Attributes.ISRC,
		})
	}

	return &songs, nil
}

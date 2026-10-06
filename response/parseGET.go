package response

import (
	"encoding/json"
	"fmt"

	"Test_api/models"
)

func ParsedResponseGet(raw []byte) ([]models.SongGet, error) {
	var respY models.RespYandexGet

	if err := json.Unmarshal(raw, &respY); err != nil {
		fmt.Printf("Unmarshal: %v\nRAW (first 500): %.500s", err, string(raw))
		return nil, fmt.Errorf("Unmarshal error: %w", err)
	}

	// fmt.Printf("Полученные данные в Get: %v\n", len(respY.Result.Tracks))
	testLen := len(respY.Result.Tracks) // len for test req and resp

	songs := make([]models.SongGet, 0, testLen) // ПОДУМАТЬ ПРО РАЗМЕР

	for id, song := range respY.Result.Tracks {
		songs = append(songs, models.SongGet{ID: song.ID, AlbumID: song.AlbumID})

		if id > 5 {
			break
		}
	}

	return songs, nil
}

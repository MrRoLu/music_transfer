package storage

import (
	"Test_api/pkg/models"
	"encoding/json"
	"os"
	"path/filepath"
)

const DefaultPath = "data/songs.json"

func SaveSongs(songs *[]models.SongYa, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")

	return encoder.Encode(songs)
}

func LoadSongs(path string) (*[]models.SongYa, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var songs []models.SongYa

	if err := json.NewDecoder(file).Decode(&songs); err != nil {
		return nil, err
	}

	return &songs, nil
}

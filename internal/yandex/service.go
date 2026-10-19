package yandex

import (
	"Test_api/internal/yandex/response"
	"Test_api/pkg/httpclient"
	"Test_api/pkg/models"

	// "Test_api/pkg/models"
	"fmt"
	"strconv"
)

type Service struct {
	client *httpclient.Client
}

func NewService() *Service {
	return &Service{
		client: newClient(),
	}
}

const (
	pathGet  = "/playlist"
	pathPost = "/tracks"
	pathLK   = "/lk.d9ec9351-ed2c-4bdd-a68d-73680bdd22bc"
)

func (s *Service) GetPlaylist(query map[string]string) (map[string][]string, error) {
	b, err := s.client.DoRequest("GET", pathGet+pathLK, httpclient.WithQuery(query))

	if err != nil {
		fmt.Printf("Error service yandex in GetPlaylist: %v\n", err)
		return nil, err
	}

	rawSongs, err := response.ParsedResponseGet(b)

	if err != nil {
		fmt.Printf("Error service yandex in ParsedResponseGet: %v\n", err)
		return nil, err
	}

	songsForPost := make(map[string][]string)

	for _, v := range rawSongs {
		strID := uintToString(v.ID)
		strAlbumID := uintToString(v.AlbumID)

		songsForPost["trackIds"] = append(songsForPost["trackIds"], strID+":"+strAlbumID)
	}

	return songsForPost, nil
}

func uintToString(v *uint) string {
	if v == nil {
		return ""
	}

	return strconv.FormatUint(uint64(*v), 10)
}

func (s *Service) GetFinallist(songsForPost map[string][]string) (*[]models.SongYa, error) {
	b, err := s.client.DoRequest("POST", pathPost, httpclient.WithMultipart(songsForPost))

	if err != nil {
		fmt.Printf("Err in service GetFinalList (POST): %v", err)
		return nil, err
	}

	resultSongs, err := response.ParsedResponsePost(b)

	if err != nil {
		fmt.Printf("Err in service ParsedResponsePost: %v", err)
		return nil, err
	}

	songs := make([]models.SongYa, 0, len(resultSongs.Result))

	for _, v := range resultSongs.Result {
		songs = append(songs, models.SongYa{
			Title:      v.Title,
			ArtistName: v.Artists[0].Name,
		})
	}

	return &songs, nil
}

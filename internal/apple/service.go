package apple

import (
	"Test_api/internal/apple/response"
	"Test_api/pkg/httpclient"
	"Test_api/pkg/models"
	"fmt"
	"log"
	"strings"
	"time"
)

type Service struct {
	client *httpclient.Client
}

func NewService(DevToken, MediaUserToken string) *Service {
	return &Service{
		client: newClient(DevToken, MediaUserToken),
	}
}

const pathStorefront = "/me/storefront"

func (s *Service) GetStorefront() (*string, error) {
	b, err := s.client.DoRequest("GET", pathStorefront)

	if err != nil {
		return nil, fmt.Errorf("get storefront: %w", err)
	}

	storefront, err := response.ParsedResponseStoreFront(b)

	if err != nil {
		return nil, err
	}

	return &storefront, nil
}

func (s *Service) FindSongs(songs *[]models.SongYa, storefront string, limit uint8) (*[]string, *[]models.SongYa) {
	pathSearhSongs := "/catalog/" + storefront + "/search"

	if limit <= 0 {
		limit = 5
	}

	songsId := make([]string, 0, len(*songs))
	var songsNotFoundId []models.SongYa

	const flush = 20

	for _, v := range *songs {
		b, err := s.client.DoRequest("GET", pathSearhSongs, httpclient.WithQuery(map[string]string{
			"term":  *v.Title + " " + *v.ArtistName,
			"types": "songs",
			"limit": fmt.Sprintf("%d", limit),
		}))

		if err != nil {
			songsNotFoundId = append(songsNotFoundId, v)
			// fmt.Printf("Error in Found songs apple: %v", err)
			continue
		}

		songsApple, err := response.ParsedCatalogSearchResult(b)

		// fmt.Printf("Found songs apple: %v", len(*songsApple))

		if err != nil {
			fmt.Printf("Error in parse ParsedCatalogSearchResult: %v", err)
			continue
		}

		right_song := pickBestMatch(&v, songsApple)

		if right_song == nil {
			songsNotFoundId = append(songsNotFoundId, v)
			continue
		}

		songsId = append(songsId, (*right_song).ID)

		if len(songsId) >= flush {
			err := s.AddSongs(&songsId)

			if err != nil {
				log.Printf("Error to add batch in library: %v\n", err)
			} else {
				log.Printf("To add in library compleate, score: %d\n", len(songsId))
				time.Sleep(1 * time.Minute)
			}

			songsId = songsId[:0]
		}
	}

	return &songsId, &songsNotFoundId

}

const pathAddSongs = "/me/library"
const batchSize = 25

func (s *Service) AddSongs(songsId *[]string) error {
	for i := 0; i < len(*songsId); i += batchSize {
		end := i + batchSize

		if end > len(*songsId) {
			end = len(*songsId)
		}

		batch := (*songsId)[i:end]

		query := map[string]string{}

		query["ids[songs]"] = strings.Join(batch, ",")

		_, err := s.client.DoRequest("POST", pathAddSongs, httpclient.WithQuery(query))

		if err != nil {
			return fmt.Errorf("Error in AddSongs, add songs batcha %v-%v, err: %w", i, end, err)
		}
	}

	return nil
}

func pickBestMatch(source *models.SongYa, candidates *[]models.Song) *models.Song {
	if len(*candidates) == 0 {
		return nil
	}

	for i, c := range *candidates {
		if strings.EqualFold(c.Name, *source.Title) && strings.EqualFold(c.ArtistName, *source.ArtistName) {
			return &(*candidates)[i]
		}
	}

	return &(*candidates)[0]
}

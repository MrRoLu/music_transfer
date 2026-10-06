package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"Test_api/request"
	"Test_api/response"
)

const (
	baseURL  = "https://api.music.yandex.ru"
	pathGet  = "/playlist"
	pathPost = "/tracks"
	pathLK   = "/lk.d9ec9351-ed2c-4bdd-a68d-73680bdd22bc"
)

// query="?resumeStream=false&richTracks=false"
func main() {
	var q = url.Values{}

	q.Set("resumeStream", "false")
	q.Set("richTracks", "false")

	respGet, err := request.DoRequestJSON(baseURL, http.MethodGet, pathGet+pathLK, q, nil)

	if err != nil {
		log.Printf("Error in doRequest: %v", err)
		return
	}

	songs, err := response.ParsedResponseGet(respGet)

	if err != nil {
		log.Printf("Error in parsed from Get func: %v", err)
		return
	}

	// for id, song := range songs {
	// 	// b, _ := json.MarshalIndent(song, "", " ")
	// 	fmt.Printf("ID: %d, SongID: %v, AlbumID: %v\n", id, *song.ID, *song.AlbumID)
	// }

	songForPost := make(map[string][]string, 10) //len(songs)

	for _, song := range songs {
		strID := uintToString(song.ID)
		strAlbumID := uintToString(song.AlbumID)

		res := strID + ":" + strAlbumID
		songForPost["trackIds"] = append(songForPost["trackIds"], res)
	}

	// fmt.Printf("Data in map after GET: %v\n", songForPost)

	// POST для Multipart
	respPost, err := request.DoRequestMultipart(baseURL, http.MethodPost, pathPost, songForPost)

	// POST для from-data (x-www-form-urlencoded)
	// respPost, err := request.DoRequestFROM(baseURL, http.MethodPost, pathPost, songForPost)

	if err != nil {
		log.Printf("Error after doRequestMultipart: %v", err)
		return
	}

	// fmt.Printf("Body in respPost: %v", string(respPost))

	resSong, err := response.ParsedResponsePost(respPost)

	if err != nil {
		log.Printf("Error in parse response post: %v", err)
		return
	}

	for id, song := range resSong.Result {
		fmt.Printf("ID элемента: %v, Song_ID: %v, Song_Title: %v, Name_artist: %v\n", id+1, *song.ID, *song.Title, *song.Artists[0].Name)

		if id == 4 {
			break
		}
	}
}

func uintToString(v *uint) string {
	if v == nil {
		return ""
	}

	return strconv.FormatUint(uint64(*v), 10)
}

// X POST\
// -H "Content-Type: application/json"\
// -H "Origin: https://music.yandex.ru"\
// -d '{"trackIds":["152866774:42729401"]}'\
// "https://api.music.yandex.ru/tracks"

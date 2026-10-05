package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	baseURL  = "https://api.music.yandex.ru"
	pathGet  = "/playlist"
	pathPost = "/tracks"
	pathLK   = "/lk.d9ec9351-ed2c-4bdd-a68d-73680bdd22bc"
)

type respYandexGet struct {
	Result struct {
		Tracks []SongGet `json:"tracks"`
	} `json:"result"`
}

type SongGet struct {
	ID      *uint `json:"id"`
	AlbumID *uint `json:"albumId"`
}

type respYandexPost struct {
	Result []struct {
		ID      *string  `json:"id"`
		RealId  *string  `json:"realId"`
		Title   *string  `json:"title"`
		Artists []artist `json:"artists"`
	} `json:"result"`
}

type artist struct {
	ID   *uint   `json:"id"`
	Name *string `json:"name"`
}

// query="?resumeStream=false&richTracks=false"
func main() {
	var q = url.Values{}

	q.Set("resumeStream", "false")
	q.Set("richTracks", "false")

	respGet, err := doRequestJSON(http.MethodGet, pathGet+pathLK, q, nil)

	if err != nil {
		log.Printf("Error in doRequest: %v", err)
		return
	}

	songs, err := parsedResponseGet(respGet)

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

	respPost, err := doRequestMultipart(http.MethodPost, pathPost, songForPost)

	if err != nil {
		log.Printf("Error after doRequestMultipart: %v", err)
		return
	}

	// fmt.Printf("Body in respPost: %v", string(respPost))

	resSong, err := parsedResponsePost(respPost)

	if err != nil {
		log.Printf("Error in parse response post: %v", err)
		return
	}

	for id, song := range resSong.Result {
		fmt.Printf("ID элемента: %v, Song_ID: %v, Song_Title: %v, Name_artist: %v\n", id+1, *song.ID, *song.Title, *song.Artists[0].Name)

		if id > 5 {
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

func parsedResponseGet(raw []byte) ([]SongGet, error) {
	var respY respYandexGet

	if err := json.Unmarshal(raw, &respY); err != nil {
		fmt.Printf("Unmarshal: %v\nRAW (first 500): %.500s", err, string(raw))
		return nil, fmt.Errorf("Unmarshal error: %w", err)
	}

	// fmt.Printf("Полученные данные в Get: %v\n", len(respY.Result.Tracks))
	testLen := len(respY.Result.Tracks) / 10 // len for test req and resp

	songs := make([]SongGet, 0, testLen)

	for id, song := range respY.Result.Tracks {
		songs = append(songs, SongGet{ID: song.ID, AlbumID: song.AlbumID})

		if id > 5 {
			break
		}
	}

	return songs, nil
}

func parsedResponsePost(raw []byte) (*respYandexPost, error) {
	var respY respYandexPost

	if err := json.Unmarshal(raw, &respY); err != nil {
		fmt.Printf("Unmarshal: %v\nRAW: %.500s\n", err, string(raw))
		return nil, fmt.Errorf("Unmarshal error: %w", err)
	}

	fmt.Printf("Полученные данные: %v\n", len(respY.Result))
	// fmt.Printf("Распарсенный Post: %v", respY.Result)
	// testLen := len(respY.Result) / 10 // len for test req and resp

	return &respY, nil
}

func doRequestJSON(method, path string, query url.Values, body any) ([]byte, error) { // переписать под JSON и под x-www-fromencode
	var reqBodyBites []byte

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			log.Printf("Marshal req body: %v", err)
			return nil, fmt.Errorf("Err in marshal: %w", err)
		}
		reqBodyBites = b
	}

	fullURL := baseURL + path

	if query != nil {
		fullURL += "?" + query.Encode()
	}

	// fmt.Printf("Full URL: %v", fullURL)

	var reqBody io.Reader
	if reqBodyBites != nil {
		reqBody = bytes.NewReader(reqBodyBites)
	}

	req, err := http.NewRequest(method, fullURL, reqBody)

	if err != nil {
		log.Printf("Request to yandex: %v", err)
		return nil, fmt.Errorf("Req to yand: %w", err)
	}

	req.Header.Set("Origin", "https://music.yandex.ru")
	req.Header.Set("Content-Type", "application/json")

	// if method == http.MethodPost {
	// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// }

	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Err in resp: %v", err)
		return nil, fmt.Errorf("Req: %w", err)
	}

	readResp, readErr := io.ReadAll(resp.Body)
	resp.Body.Close()

	if readErr != nil {
		log.Printf("Err in readResp: %v", readErr)
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("yandex api error: status=%d body=%s", resp.StatusCode, string(readResp))
	}

	return readResp, nil
}

func doRequestMultipart(method, path string, data map[string][]string) ([]byte, error) { // переписать под Multipart
	var body bytes.Buffer
	w := multipart.NewWriter(&body)

	for k, values := range data {
		for _, v := range values {
			if err := w.WriteField(k, v); err != nil {
				return nil, fmt.Errorf("Error to WriteField: %w", err)
			}
		}
	}

	w.Close()

	fullURL := baseURL + path

	// fmt.Printf("Full URL: %v", fullURL)

	// fmt.Printf("Body in reqMultipart: %v\n", body.String())

	req, err := http.NewRequest(method, fullURL, &body)

	if err != nil {
		log.Printf("Request to yandex: %v", err)
		return nil, fmt.Errorf("Req to yand: %w", err)
	}

	req.Header.Set("Origin", "https://music.yandex.ru")
	req.Header.Set("Content-Type", w.FormDataContentType())

	// if method == http.MethodPost {
	// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// }

	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Err in resp: %v", err)
		return nil, fmt.Errorf("Req: %w", err)
	}

	readResp, readErr := io.ReadAll(resp.Body)
	resp.Body.Close()

	if readErr != nil {
		log.Printf("Err in readResp: %v", readErr)
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("yandex api error: status=%d body=%s", resp.StatusCode, string(readResp))
	}

	return readResp, nil
}

// func doRequestFROM(method, path string, query url.Values, data []string) ([]byte, error) { // переписать под x-www-fromencode
// 	var body url.Values

// 	if data != nil {
// 		for _, v := range data {
// 			body.Add("track-ids", v)
// 		}
// 	}

// 	body.Add("")

// 	// if body != nil {
// 	// 	b, err := json.Marshal(body)
// 	// 	if err != nil {
// 	// 		log.Printf("Marshal req body: %v", err)
// 	// 		return nil, fmt.Errorf("Err in marshal: %w", err)
// 	// 	}
// 	// 	reqBodyBites = b
// 	// }

// 	fullURL := baseURL + path

// 	if query != nil {
// 		fullURL += "?" + query.Encode()
// 	}

// 	// fmt.Printf("Full URL: %v", fullURL)

// 	// var reqBody io.Reader
// 	// if reqBodyBites != nil {
// 	// 	reqBody = bytes.NewReader(reqBodyBites)
// 	// }

// 	req, err := http.NewRequest(method, fullURL, reqBody)

// 	if err != nil {
// 		log.Printf("Request to yandex: %v", err)
// 		return nil, fmt.Errorf("Req to yand: %w", err)
// 	}

// 	req.Header.Set("Origin", "https://music.yandex.ru")
// 	req.Header.Set("Content-Type", "application/json")

// 	// if method == http.MethodPost {
// 	// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 	// }

// 	client := &http.Client{Timeout: 30 * time.Second}

// 	resp, err := client.Do(req)

// 	if err != nil {
// 		log.Printf("Err in resp: %v", err)
// 		return nil, fmt.Errorf("Req: %w", err)
// 	}

// 	readResp, readErr := io.ReadAll(resp.Body)
// 	resp.Body.Close()

// 	if readErr != nil {
// 		log.Printf("Err in readResp: %v", readErr)
// 	}

// 	if resp.StatusCode >= 300 {
// 		return nil, fmt.Errorf("yandex api error: status=%d body=%s", resp.StatusCode, string(readResp))
// 	}

// 	return readResp, nil
// }

// X POST\
// -H "Content-Type: application/json"\
// -H "Origin: https://music.yandex.ru"\
// -d '{"trackIds":["152866774:42729401"]}'\
// "https://api.music.yandex.ru/tracks"

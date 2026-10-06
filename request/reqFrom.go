package request

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func DoRequestFROM(baseURL, method, path string, data map[string][]string) ([]byte, error) { // x-www-fromencode ПРОТЕСТИРОВАТЬ
	var body url.Values

	for key, value := range data {
		for _, v := range value {
			body.Add(key, v)
		}
	}

	// body.Add("") // добавить доп параметр в тело

	fullURL := baseURL + path

	// fmt.Printf("Full URL: %v", fullURL)

	// var reqBody io.Reader
	// if reqBodyBites != nil {
	// 	reqBody = bytes.NewReader(reqBodyBites)
	// }

	req, err := http.NewRequest(method, fullURL, strings.NewReader(body.Encode()))

	if err != nil {
		log.Printf("Request to yandex: %v", err)
		return nil, fmt.Errorf("Req to yand: %w", err)
	}

	req.Header.Set("Origin", "https://music.yandex.ru")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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

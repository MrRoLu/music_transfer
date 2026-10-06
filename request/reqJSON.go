package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func DoRequestJSON(baseURL, method, path string, query url.Values, body any) ([]byte, error) {
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

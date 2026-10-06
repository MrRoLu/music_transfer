package request

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

func DoRequestMultipart(baseURL, method, path string, data map[string][]string) ([]byte, error) {
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

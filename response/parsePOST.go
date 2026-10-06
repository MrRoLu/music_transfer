package response

import (
	"Test_api/models"
	"encoding/json"
	"fmt"
)

func ParsedResponsePost(raw []byte) (*models.RespYandexPost, error) {
	var respY models.RespYandexPost

	if err := json.Unmarshal(raw, &respY); err != nil {
		fmt.Printf("Unmarshal: %v\nRAW: %.500s\n", err, string(raw))
		return nil, fmt.Errorf("Unmarshal error: %w", err)
	}

	fmt.Printf("Полученные данные: %v\n", len(respY.Result))
	// fmt.Printf("Распарсенный Post: %v", respY.Result)
	// testLen := len(respY.Result) / 10 // len for test req and resp

	return &respY, nil
}

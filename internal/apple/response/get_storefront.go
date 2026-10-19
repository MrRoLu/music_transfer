package response

import (
	"Test_api/pkg/models"
	"encoding/json"
	"fmt"
)

func ParsedResponseStoreFront(raw []byte) (string, error) {
	var parsed models.StorefrontResponse

	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", err
	}
	if len(parsed.Data) == 0 {
		return "", fmt.Errorf("get storefront: empty response data")
	}

	return parsed.Data[0].ID, nil
}

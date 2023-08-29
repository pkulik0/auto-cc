package deepl

import "encoding/json"

type UsageInfo struct {
	CharactersUsed  int64 `json:"character_count"`
	CharactersLimit int64 `json:"character_limit"`
}

func (d *DeepL) GetUsage() (*UsageInfo, error) {
	resp, err := d.request("/usage", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var usageInfo UsageInfo
	if err := json.NewDecoder(resp.Body).Decode(&usageInfo); err != nil {
		return nil, err
	}

	return &usageInfo, nil
}

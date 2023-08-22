package deepl

import (
	"encoding/json"
	"io"
	"net/http"
)

type DeepL struct {
	apiKey  string
	baseUrl string
}

func NewDeepL(baseUrl string, apiKey string) *DeepL {
	return &DeepL{
		apiKey:  apiKey,
		baseUrl: baseUrl,
	}
}

func (d *DeepL) request(endpoint string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("GET", d.baseUrl+endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "DeepL-Auth-Key "+d.apiKey)

	client := http.Client{}
	return client.Do(req)
}

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

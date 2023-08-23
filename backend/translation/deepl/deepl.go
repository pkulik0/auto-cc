package deepl

import (
	"encoding/json"
	"errors"
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
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "DeepL-Auth-Key "+d.apiKey)

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		var errorResponse struct {
			Message string
		}
		if err := json.NewDecoder(res.Body).Decode(&errorResponse); err != nil {
			return nil, err
		}
		return nil, errors.New("deepl api error: " + errorResponse.Message)
	}

	return res, nil
}

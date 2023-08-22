package deepl

import (
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

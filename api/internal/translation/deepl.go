package translation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkulik0/autocc/api/internal/store"
)

const (
	baseURL = "https://api-free.deepl.com/v2/"
)

type deeplTransport struct {
	base   http.RoundTripper
	apiKey string
}

func newDeeplTransport(base http.RoundTripper, apiKey string) *deeplTransport {
	return &deeplTransport{
		base:   base,
		apiKey: apiKey,
	}
}

func (t *deeplTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "DeepL-Auth-Key "+t.apiKey)
	return t.base.RoundTrip(req)
}

type deeplApiClient struct {
	client     *http.Client
	revertCost func() error
}

func newDeeplApiClient(ctx context.Context, store store.Store, neededQuota uint) (*deeplApiClient, error) {
	credentials, revert, err := store.GetCredentialsDeepLByAvailableCost(ctx, neededQuota)
	if err != nil {
		return nil, err
	}

	return &deeplApiClient{
		client: &http.Client{
			Transport: newDeeplTransport(http.DefaultTransport, credentials.Key),
		},
		revertCost: revert,
	}, nil
}

func (c *deeplApiClient) request(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, baseURL+path, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.client.Do(req)
}

type languageResponse struct {
	Language string `json:"language"`
	Name     string `json:"name"`
}

func (c *deeplApiClient) getLanguages(ctx context.Context) ([]string, error) {
	resp, err := c.request(ctx, http.MethodGet, "languages", nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data []languageResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	languages := make([]string, len(data))
	for i, l := range data {
		languages[i] = l.Language
	}

	return languages, nil
}

type translateRequest struct {
	Text           []string `json:"text"`
	SourceLanguage string   `json:"source_lang"`
	TargetLanguage string   `json:"target_lang"`
}

type translateResponse struct {
	Translations []struct {
		DetectedSourceLanguage string `json:"detected_source_language"`
		Text                   string `json:"text"`
	} `json:"translations"`
}

func (c *deeplApiClient) translate(ctx context.Context, text []string, sourceLanguage, targetLanguage string) ([]string, error) {
	data, err := json.Marshal(&translateRequest{
		Text:           text,
		SourceLanguage: sourceLanguage,
		TargetLanguage: targetLanguage,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.request(ctx, http.MethodPost, "translate", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
	}

	result := &translateResponse{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	translations := make([]string, len(result.Translations))
	for i, t := range result.Translations {
		translations[i] = t.Text
	}

	return translations, nil
}

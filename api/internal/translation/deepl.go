package translation

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkulik0/autocc/api/internal/store"
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
	baseUrl    string
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
		baseUrl:    "https://api.deepl.com/v2/",
	}, nil
}

type languagesResponse struct {
	Data []struct {
		Language          string `json:"language"`
		Name              string `json:"name"`
		SupportsFormality bool   `json:"supportsFormality"`
	}
}

func (r *languagesResponse) toLanguages() []Language {
	languages := make([]Language, 0, len(r.Data))
	for _, l := range r.Data {
		languages = append(languages, Language(l.Language))
	}
	return languages
}

func (c *deeplApiClient) request(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseUrl+path, body)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

func (c *deeplApiClient) getLanguages(ctx context.Context) ([]Language, error) {
	resp, err := c.request(ctx, http.MethodGet, "languages", nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var data languagesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.toLanguages(), nil
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
	}
}

func (c *deeplApiClient) translate(ctx context.Context, text []string, source, target Language) ([]string, error) {
	data, err := json.Marshal(translateRequest{
		Text:           text,
		SourceLanguage: source.String(),
		TargetLanguage: target.String(),
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.request(ctx, http.MethodPost, "translate", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var result translateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	translations := make([]string, 0, len(result.Translations))
	for i, t := range result.Translations {
		translations[i] = t.Text
	}

	return translations, nil
}

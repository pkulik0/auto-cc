package deepl

import (
	"bytes"
	"encoding/json"
)

type translateRequest struct {
	Text                     []string `json:"text"`
	SourceLangCode           string   `json:"source_lang"`
	TargetLangCode           string   `json:"target_lang"`
	ShouldPreserveFormatting bool     `json:"preserve_formatting"`
}

type translateResponse struct {
	Translations []struct {
		Text string
	}
}

func newTranslateRequest(text []string, sourceLangCode string, targetLangCode string) translateRequest {
	return translateRequest{
		Text:                     text,
		SourceLangCode:           sourceLangCode,
		TargetLangCode:           targetLangCode,
		ShouldPreserveFormatting: true,
	}
}

func (d *DeepL) Translate(text []string, sourceLangCode string, targetLangCode string) ([]string, error) {
	requestData := newTranslateRequest(text, sourceLangCode, targetLangCode)
	requestJson, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	resp, err := d.request("/translate", bytes.NewBuffer(requestJson))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var parsedResponse translateResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsedResponse); err != nil {
		return nil, err
	}

	translatedText := make([]string, 0, len(parsedResponse.Translations))
	for _, translation := range parsedResponse.Translations {
		translatedText = append(translatedText, translation.Text)
	}

	return translatedText, nil
}

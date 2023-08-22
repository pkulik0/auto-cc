package deepl

import (
	"encoding/json"
)

type Language struct {
	Code string `json:"language"`
	Name string `json:"name"`
}

type LanguageType string

const (
	SourceLanguages LanguageType = "source"
	TargetLanguages LanguageType = "target"
)

func (d *DeepL) GetLanguages(languageType LanguageType) ([]Language, error) {
	resp, err := d.request("/languages?type="+string(languageType), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var languages []Language
	if err := json.NewDecoder(resp.Body).Decode(&languages); err != nil {
		return nil, err
	}

	return languages, nil
}

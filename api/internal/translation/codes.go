package translation

import "strings"

// CodeGoogleToTranslation translates the language code from the Google API format.
func CodeGoogleToTranslation(code string) string {
	switch strings.ToLower(code) {
	case "no":
		return "nb"
	default:
		return code
	}
}

// CodeTranslationToGoogle translates the language code to the Google API format.
func CodeTranslationToGoogle(code string) string {
	switch strings.ToLower(code) {
	case "nb":
		return "no"
	default:
		return code
	}
}

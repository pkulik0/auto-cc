import type {Language} from "$lib/languages/api";

export let targetLanguages: string[] = []

export const setTargetLanguages = (languages: Language[]) => {
    targetLanguages = [...new Set(languages.map(language => language.language.split("-")[0]))]
}

export let sourceLanguageCode = ""
export const setSourceLanguageCode = (code: string) => {
    sourceLanguageCode = code
}

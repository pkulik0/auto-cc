import {PUBLIC_LANG_URL} from "$env/static/public";

export interface Language {
    name: string
    language: string
}

interface LanguagesResponse {
    source: Language[],
    target: Language[]
}

export const getLanguages = async (): Promise<LanguagesResponse> => {
    const response = await fetch(PUBLIC_LANG_URL+"/languages")
    if(!response.ok) {
        throw new Error("Failed to fetch available languages.")
    }
    return response.json()
}

export const translateText = async (text: string[], sourceLanguageCode: string, targetLanguageCode: string): Promise<string[]> => {
    const response = await fetch(PUBLIC_LANG_URL+"/translate", {
        method: "POST",
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            "text": text,
            "source": sourceLanguageCode,
            "target": targetLanguageCode,
        }),
    })
    if(!response.ok) {
        console.log(await response.text())
        throw new Error("Failed to translate text")
    }
    return response.json()
}
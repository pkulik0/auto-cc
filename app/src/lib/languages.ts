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
import type {Language} from "$lib/languages/api";
import {writable} from "svelte/store";

export let targetLanguages: string[] = []
export let filteredTargetLanguages: string[] = []

export const setTargetLanguages = (languages: Language[]) => {
    targetLanguages = [...new Set(languages.map(language => language.language.split("-")[0]))]
}

export let sourceLanguageCode = ""
export const setSourceLanguageCode = (code: string) => {
    sourceLanguageCode = code
}

const savedLanguageKey = "selectedLanguage"
const getSavedLanguage = (): Language|null => {
    if(typeof window === "undefined") return null
    const entry = window.localStorage.getItem(savedLanguageKey)
    if(!entry) return null
    return JSON.parse(entry)
}

export const selectedLanguage = writable<Language|null>(getSavedLanguage())
selectedLanguage.subscribe(value => {
    if(!value) return

    filteredTargetLanguages = targetLanguages.filter((code: string) => code !== value.language)
    setSourceLanguageCode(value.language)

    if(typeof window === "undefined") return
    window.localStorage.setItem(savedLanguageKey, JSON.stringify(value))
})
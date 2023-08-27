import type {Language} from "$lib/languages/api";
import {writable} from "svelte/store";

export let targetLanguages: string[] = []
export let filteredTargetLanguages: string[] = []

export const setTargetLanguages = (languages: Language[]) => {
    targetLanguages = [...new Set(languages.map(language => language.language.split("-")[0]))]
    filteredTargetLanguages = targetLanguages.filter((code: string) => code !== sourceLanguageCode)
}

const savedLanguageStorageKey = "selectedLanguage"
const getSavedLanguage = (): Language|null => {
    if(typeof window === "undefined") return null
    const entry = window.localStorage.getItem(savedLanguageStorageKey)
    if(!entry) return null
    return JSON.parse(entry)
}

export let sourceLanguageCode = ""
export const selectedLanguage = writable<Language|null>(getSavedLanguage())
selectedLanguage.subscribe(value => {
    if(!value) return

    sourceLanguageCode = value.language
    filteredTargetLanguages = targetLanguages.filter((code: string) => code !== sourceLanguageCode)

    if(typeof window === "undefined") return
    window.localStorage.setItem(savedLanguageStorageKey, JSON.stringify(value))
})
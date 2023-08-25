import {writable} from "svelte/store";
import type {Language} from "$lib/languages/api";
import {setSourceLanguageCode} from "$lib/languages/languages";

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
    setSourceLanguageCode(value.language)
    if(typeof window === "undefined") return
    window.localStorage.setItem(savedLanguageKey, JSON.stringify(value))
})
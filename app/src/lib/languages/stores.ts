import {writable} from "svelte/store";
import type {Language} from "$lib/languages/api";

const savedLanguageKey = "selectedLanguage"
const getSavedLanguage = (): Language|null => {
    if(typeof window === "undefined") return null
    const entry = window.localStorage.getItem(savedLanguageKey)
    if(!entry) return null
    return JSON.parse(entry)
}

export const selectedLanguage = writable<Language|null>(getSavedLanguage())
selectedLanguage.subscribe(value => {
    if(typeof window === "undefined") return
    window.localStorage.setItem(savedLanguageKey, JSON.stringify(value))
})

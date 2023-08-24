import axios from "axios";
import {PUBLIC_API_URL} from "$env/static/public";

export interface Language {
    name: string
    language: string
}

interface LanguagesResponse {
    source: Language[],
    target: Language[]
}

export const getLanguages = async (): Promise<LanguagesResponse> => {
    const response = await axios.get<LanguagesResponse>(PUBLIC_API_URL+"/languages")
    if(response.status != 200) {
        throw new Error("Couldn't fetch available languages.")
    }
    return response.data
}
import {translateText} from "$lib/languages/api";
import {targetLanguages} from "$lib/languages/data";
import {env} from "$env/dynamic/public";

export interface VideoMetadata {
    title: string
    description: string
    language: string
}

export const metadataSeparator = ";"

export const getMetadata = async (videoId: string): Promise<VideoMetadata> => {
    const response = await fetch(env.PUBLIC_API_URL+`/youtube/videos/${videoId}`)
    if(!response.ok) throw new Error(`Failed to fetch metadata for ${videoId}`)
    return response.json()
}

export const translateMetadata = async (metadata: VideoMetadata): Promise<VideoMetadata[]> => {
    const splitTitle = metadata.title.split(metadataSeparator)
    const splitDescription = metadata.description.split(metadataSeparator)
    const text: string[] = [...splitTitle, ...splitDescription]

    const promises = targetLanguages.map(targetLanguageCode => translateText(text, metadata.language, targetLanguageCode))
    const translatedTexts = await Promise.all(promises)

    return translatedTexts.map((translatedText, index) => {
        if(translatedText.length < 2) throw new Error("Not enough texts received from translation")

        const title = translatedText.slice(0, splitTitle.length).join("")
        const description = translatedText.slice(splitTitle.length, translatedText.length).join("")
        const language = targetLanguages[index]

        return {
            "title": title,
            "description": description,
            "language": language
        }
    })
}

export const insertVideoMetadata = async (videoId: string, metadataArray: VideoMetadata[]) => {
    const response = await fetch(env.PUBLIC_API_URL+`/youtube/videos/${videoId}`, {
        method: "POST",
        body: JSON.stringify(metadataArray),
        headers: {
            "Content-Type": "application/json"
        }
    })
    if(!response.ok) throw new Error("Failed to update video metadata")
}
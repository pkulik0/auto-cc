import type {Video} from "$lib/youtube/video";
import {Srt, translateSrt} from "$lib/youtube/srt";
import type {Language} from "$lib/languages/api";
import {PUBLIC_YOUTUBE_URL} from "$env/static/public";

export const translateVideoCC = async (sourceLanguage: Language, video: Video) => {
    const languageCode: string = sourceLanguage.language.toLowerCase()

    const ccList: CCEntry[] = await getCCList(video.id)
    if(!ccList) {
        throw new Error("No CCs found.")
    }

    const ccEntry = ccList.find(cc => cc.language === languageCode)
    if(!ccEntry) {
        throw new Error(`${sourceLanguage.name} CC not found!`)
    }

    const srt = new Srt(await downloadCC(ccEntry.id))
    const translatedSrts = await translateSrt(srt)
    console.log(translatedSrts)
}

interface CCEntry {
    id: string,
    language: string
}

const getCCList = async (videoId: string): Promise<CCEntry[]> => {
    const response = await fetch(PUBLIC_YOUTUBE_URL+`/videos/${videoId}/cc`)
    if(!response.ok) {
        throw new Error(`Failed to fetch CCs of ${videoId}`)
    }
    return response.json()
}

const downloadCC = async (ccId: string): Promise<string> => {
    const response = await fetch(PUBLIC_YOUTUBE_URL+`/cc/${ccId}`)
    if(!response.ok) {
        throw new Error(`Failed to download CC with id ${ccId}`)
    }
    return (await response.text()).trimEnd()
}
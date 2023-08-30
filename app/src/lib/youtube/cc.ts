import type {Video} from "$lib/youtube/video";
import {Srt, translateSrt} from "$lib/youtube/srt";
import {env} from "$env/dynamic/public";

export const translateVideoCC = async (video: Video, sourceLanguageCode: string, targetLanguagesCodes: string[]) => {
    const ccList: CCEntry[] = await getCCList(video.id)
    if(!ccList) throw new Error("No CCs found.")

    const ccEntry = ccList.find(cc => cc.language === sourceLanguageCode.toLowerCase())
    if(!ccEntry) throw new Error(`${sourceLanguageCode} CC not found!`)

    const srt = new Srt(await downloadCC(ccEntry.id))
    const translatedSrts = await translateSrt(srt, sourceLanguageCode, targetLanguagesCodes)

    for(const [index, translatedSrt] of translatedSrts.entries()) {
        await insertCC(srt, targetLanguagesCodes[index], video.id)
    }
}

interface CCEntry {
    id: string,
    language: string
}

const getCCList = async (videoId: string): Promise<CCEntry[]> => {
    const response = await fetch(env.PUBLIC_API_URL+`/youtube/videos/${videoId}/cc`)
    if(!response.ok) throw new Error(`Failed to fetch CCs of ${videoId}`)
    return response.json()
}

const downloadCC = async (ccId: string): Promise<string> => {
    const response = await fetch(env.PUBLIC_API_URL+`/youtube/cc/${ccId}`)
    if(!response.ok) throw new Error(`Failed to download CC with id ${ccId}`)
    return (await response.text()).trimEnd()
}

const insertCC = async (srt: Srt, languageCode: string, videoId: string): Promise<void> => {
    const response = await fetch(env.PUBLIC_API_URL+`/youtube/videos/${videoId}/cc?language=${languageCode}`, {
        method: "POST",
        body: srt.toString(),
    })
    if(!response.ok) throw new Error("Failed to insert CC")
}
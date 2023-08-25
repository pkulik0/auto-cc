import {PUBLIC_YOUTUBE_URL} from "$env/static/public";

export interface Video {
    id: string,
    title: string,
    thumbnailUrl: string,
    description: string,
    publishedAt: string
}

export const getVideos = async (): Promise<Video[]> => {
    const response = await fetch(PUBLIC_YOUTUBE_URL+"/videos")
    if(!response.ok) {
        throw new Error(`Failed to fetch videos`)
    }
    return response.json()
}

export interface CCEntry {
    id: string,
    language: string
}

export const getCCList = async (videoId: string): Promise<CCEntry[]> => {
    const response = await fetch(PUBLIC_YOUTUBE_URL+`/videos/${videoId}/cc`)
    if(!response.ok) {
        throw new Error(`Failed to fetch CCs of ${videoId}`)
    }
    return response.json()
}

export const downloadCC = async (ccId: string): Promise<string> => {
    const response = await fetch(PUBLIC_YOUTUBE_URL+`/cc/${ccId}`)
    if(!response.ok) {
        throw new Error(`Failed to download CC with id ${ccId}`)
    }
    return response.text()
}
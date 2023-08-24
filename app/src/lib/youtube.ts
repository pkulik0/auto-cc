import {PUBLIC_YOUTUBE_URL} from "$env/static/public";

export interface Video {
    id: string,
    title: string,
    thumbnailUrl: string,
    description: string,
    publishedAt: string
}

const get = async (endpoint: string): Promise<any> => {
    const response = await fetch(PUBLIC_YOUTUBE_URL+endpoint)
    if(!response.ok) {
        throw new Error(`Failed to GET: ${PUBLIC_YOUTUBE_URL}${endpoint}`)
    }
    return response.json()
}

export const getVideos = async (): Promise<Video[]> => {
    return get("/videos")
}
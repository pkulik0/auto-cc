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
        throw new Error("Failed to fetch videos.")
    }
    return response.json()
}
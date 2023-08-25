import {PUBLIC_YOUTUBE_URL} from "$env/static/public";
import {writable} from "svelte/store";

export interface Video {
    id: string,
    title: string,
    thumbnailUrl: string,
    description: string,
    publishedAt: string
}

export const getVideos = async (fresh = false): Promise<Video[]> => {
    const response = await fetch(PUBLIC_YOUTUBE_URL+"/videos", {
        headers: fresh ? { "Cache-Control": "no-cache" } : {}
    })
    if(!response.ok) {
        throw new Error(`Failed to fetch videos`)
    }
    return response.json()
}

export const videos = writable<Video[]>([])
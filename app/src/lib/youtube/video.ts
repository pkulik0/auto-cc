import {PUBLIC_API_URL} from "$env/static/public";
import {writable} from "svelte/store";

export interface Video {
    id: string,
    title: string,
    thumbnailUrl: string,
    description: string,
    publishedAt: string
}

interface VideosResponse {
    videos: Video[],
    nextPageToken: string
}

export const videos = writable<Video[]>([])
export const videosNextPageToken = writable<string>("")


export const getVideos = async (fresh = false, next = false): Promise<Video[]> => {
    let url = PUBLIC_API_URL+"/youtube/videos"

    let token = ""
    videosNextPageToken.subscribe(s => token = s)()
    if(next && token) url += "?token=" + token

    const response = await fetch(url, {
        headers: fresh ? { "Cache-Control": "no-cache" } : {}
    })
    if(!response.ok) {
        throw new Error(`Failed to fetch videos`)
    }

    const videosResponse: VideosResponse = await response.json()
    videosNextPageToken.set(videosResponse.nextPageToken)

    return videosResponse.videos.filter(v => !v.description.includes("#short"))
}
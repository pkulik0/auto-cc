import {PUBLIC_YOUTUBE_URL} from "$env/static/public";

export interface Identity {
    identityHash: string
    usedQuota: number
    isSelected: boolean
}

export const getIdentities = async (): Promise<Identity[]> => {
    const response = await fetch(PUBLIC_YOUTUBE_URL+"/identities")
    if(!response.ok) throw new Error("Failed to get identities")
    return response.json()
}
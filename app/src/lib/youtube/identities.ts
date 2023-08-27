import {PUBLIC_YOUTUBE_URL} from "$env/static/public";

export interface Identity {
    identityHash: string
    usedQuota: number
    isSelected: boolean
}

export const getIdentities = async (): Promise<Identity[]> => {
    const response = await fetch(PUBLIC_YOUTUBE_URL+"/identities")
    if(!response.ok) throw new Error("Failed to get identities")
    const identities: Identity[] = await response.json()
    return identities.sort((i1, i2) => i1.identityHash.localeCompare(i2.identityHash))
}

export const addIdentity = async (clientId: string, clientSecret: string) => {
    const response = await fetch(PUBLIC_YOUTUBE_URL+"/identities", {
        method: "POST",
        body: JSON.stringify({
            "clientId": clientId,
            "clientSecret": clientSecret
        }),
        headers: {
            "Content-Type": "application/json"
        }
    })
    if(!response.ok) throw new Error("Failed to add identity")
}
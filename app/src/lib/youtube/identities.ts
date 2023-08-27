import {PUBLIC_YOUTUBE_URL} from "$env/static/public";

interface Identity {
    identityHash: string
    usedQuota: number
    isSelected: boolean
}

interface AuthInfo {
    identityHash: string
    url: string
}

export interface IdentityInfo {
    hash: string,
    usedQuota: number,
    isSelected: boolean,
    authUrl: string|undefined
}

export const getIdentityInfos = async () => {
    const identities = await getIdentities()
    const authInfos = await getAuthInfo()

    const identityInfos: IdentityInfo[] = identities.map(identity => {
        const matchingInfo = authInfos.find(info => info.identityHash == identity.identityHash)
        return {
            "hash": identity.identityHash,
            "usedQuota": identity.usedQuota,
            "isSelected": identity.isSelected,
            "authUrl": matchingInfo?.url
        }
    })

    console.log(identityInfos)
    return identityInfos
}

const getAuthInfo = async (): Promise<AuthInfo[]> => {
    const response = await fetch(PUBLIC_YOUTUBE_URL+"/auth")
    if(!response.ok) throw new Error("Failed to get auth info")
    return response.json()
}

const getIdentities = async (): Promise<Identity[]> => {
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
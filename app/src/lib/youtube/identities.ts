import {env} from "$env/dynamic/public";

interface Identity {
    identityHash: string
    usedQuota: number
    isSelected: boolean
}

interface AuthUrl {
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
    const authUrls = await getAuthUrls()

    const identityInfos: IdentityInfo[] = identities.map(identity => {
        const matchingAuthUrl = authUrls.find(info => info.identityHash == identity.identityHash)
        return {
            "hash": identity.identityHash,
            "usedQuota": identity.usedQuota,
            "isSelected": identity.isSelected,
            "authUrl": matchingAuthUrl?.url
        }
    })

    return identityInfos
}

const getAuthUrls = async (): Promise<AuthUrl[]> => {
    const response = await fetch(env.PUBLIC_API_URL+"/youtube/auth")
    if(!response.ok) throw new Error("Failed to get auth info")
    return response.json()
}

const getIdentities = async (): Promise<Identity[]> => {
    const response = await fetch(env.PUBLIC_API_URL+"/youtube/identities")
    if(!response.ok) throw new Error("Failed to get identities")
    const identities: Identity[] = await response.json()
    return identities.sort((i1, i2) => i1.identityHash.localeCompare(i2.identityHash))
}

export const addIdentity = async (clientId: string, clientSecret: string) => {
    const response = await fetch(env.PUBLIC_API_URL+"/youtube/identities", {
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
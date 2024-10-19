import { PUBLIC_API_URL } from "$env/static/public"
import { userManager } from "./auth"
import { AddCredentialsDeepLRequest, AddCredentialsDeepLResponse, AddCredentialsGoogleRequest, AddCredentialsGoogleResponse, CredentialsDeepL, CredentialsGoogle, GetCredentialsResponse, GetSessionGoogleURLResponse, GetUserSessionsGoogleResponse } from "./pb/credentials"
import { GetMetadataResponse, GetYoutubeVideosResponse, Metadata, UpdateMetadataRequest } from "./pb/youtube"

const getApiUrl = (endpoint: string) => {
    if (!endpoint.startsWith("/")) endpoint = "/" + endpoint
    return PUBLIC_API_URL + endpoint
}

export interface Credentials {
    google: CredentialsGoogle[]
    deepl: CredentialsDeepL[]
}

export const getCredentials = async (): Promise<Credentials> => {
    const u = await userManager.getUser()
    if (!u) throw new Error("User not logged in")
    const token = u.access_token

    const res = await fetch(getApiUrl("/credentials"), {
        headers: {
            Authorization: `Bearer ${token}`
        }
    });
    if (!res.ok) {
        throw new Error("Failed to get clients")
    }

    const data = await res.arrayBuffer()
    const resp = GetCredentialsResponse.decode(new Uint8Array(data))
    return  { google: resp.google, deepl: resp.deepl }
}

export const addCredentialsGoogle = async (clientId: string, clientSecret: string): Promise<CredentialsGoogle> => {
    const u = await userManager.getUser()
    if (!u) throw new Error("User not logged in")
    const token = u.access_token

    const res = await fetch(getApiUrl("/credentials/google"), {
        method: "POST",
        headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/octet-stream"
        },
        body: AddCredentialsGoogleRequest.encode({
            clientId: clientId,
            clientSecret: clientSecret,
        }).finish()
    });
    if (!res.ok) {
        throw new Error("Failed to add client credential")
    }

    const data = await res.arrayBuffer()
    const credentials = AddCredentialsGoogleResponse.decode(new Uint8Array(data)).credentials

    if (!credentials) throw new Error("Failed to add client credential")
    return credentials
}

export const addCredentialsDeepL = async (key: string): Promise<CredentialsDeepL> => {
    const u = await userManager.getUser()
    if (!u) throw new Error("User not logged in")
    const token = u.access_token

    const res = await fetch(getApiUrl("/credentials/deepl"), {
        method: "POST",
        headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/octet-stream"
        },
        body: AddCredentialsDeepLRequest.encode({
            key: key,
        }).finish()
    });
    if (!res.ok) {
        throw new Error("Failed to add client credential")
    }

    const data = await res.arrayBuffer()
    const credentials = AddCredentialsDeepLResponse.decode(new Uint8Array(data)).credentials
    
    if (!credentials) throw new Error("Failed to add client credential")
    return credentials
}

export const removeCredentials = async (type: "google" | "deepl", id: number): Promise<void> => {
    const u = await userManager.getUser()
    if (!u) throw new Error("User not logged in")
    const token = u.access_token

    const res = await fetch(getApiUrl(`/credentials/${type}/${id}`), {
        method: "DELETE",
        headers: {
            Authorization: `Bearer ${token}`,
        }
    });
    if (!res.ok) {
        throw new Error("Failed to remove client credential")
    }
}

export const getUserSessionsGoogle = async (): Promise<number[]> => {
    const u = await userManager.getUser()
    if (!u) throw new Error("User not logged in")
    const token = u.access_token

    const res = await fetch(getApiUrl("/sessions/google"), {
        headers: {
            Authorization: `Bearer ${token}`
        }
    });
    if (!res.ok) {
        throw new Error("Failed to get sessions")
    }

    const data = await res.arrayBuffer()
    return GetUserSessionsGoogleResponse.decode(new Uint8Array(data)).credentialIds
}

export const getSessionGoogleURL = async (id: number): Promise<string> => {
    const u = await userManager.getUser()
    if (!u) throw new Error("User not logged in")
    const token = u.access_token

    const redirectUrl = encodeURIComponent(`${window.location.href}`)
    const res = await fetch(getApiUrl(`/sessions/google/${id}?redirect_url=${redirectUrl}`), {
        headers: {
            Authorization: `Bearer ${token}`
        }
    });
    if (!res.ok) {
        throw new Error("Failed to get session URL")
    }

    const data = await res.arrayBuffer()
    return GetSessionGoogleURLResponse.decode(new Uint8Array(data)).url
}

export const removeSessionGoogle = async (id: number): Promise<void> => {
    const u = await userManager.getUser()
    if (!u) throw new Error("User not logged in")
    const token = u.access_token

    const res = await fetch(getApiUrl(`/sessions/google/${id}`), {
        method: "DELETE",
        headers: {
            Authorization: `Bearer ${token}`,
        }
    });
    if (!res.ok) {
        throw new Error("Failed to remove session")
    }
}

export const getVideos = async (nextPageToken?: string): Promise<GetYoutubeVideosResponse> => {
    const u = await userManager.getUser()
    if (!u) throw new Error("User not logged in")
    const token = u.access_token

    const res = await fetch(getApiUrl(`/youtube/videos?next_page_token=${nextPageToken || ""}`), {
        headers: {
            Authorization: `Bearer ${token}`
        }
    });
    if (!res.ok) {
        throw new Error("Failed to get videos")
    }

    const data = await res.arrayBuffer()
    return GetYoutubeVideosResponse.decode(new Uint8Array(data))
}

export const getMetadata = async (videoId: string): Promise<GetMetadataResponse> => {
    const u = await userManager.getUser()
    if (!u) throw new Error("User not logged in")
    const token = u.access_token

    const res = await fetch(getApiUrl(`/youtube/metadata?id=${videoId}`), {
        headers: {
            Authorization: `Bearer ${token}`,
        },
    });
    if (!res.ok) {
        throw new Error("Failed to get metadata")
    }

    const data = await res.arrayBuffer()
    return GetMetadataResponse.decode(new Uint8Array(data))
}

export const updateMetadata = async (videoId: string, metadata: { [langCode: string]: Metadata }): Promise<void> => {
    const u = await userManager.getUser()
    if (!u) throw new Error("User not logged in")
    const token = u.access_token

    const res = await fetch(getApiUrl(`/youtube/metadata`), {
        method: "PUT",
        headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/octet-stream"
        },
        body: UpdateMetadataRequest.encode({ id: videoId, metadata: metadata }).finish()
    });
    if (!res.ok) {
        throw new Error("Failed to update metadata")
    }
}
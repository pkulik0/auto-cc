import { PUBLIC_API_URL } from "$env/static/public"
import { userManager } from "./auth"

const getApiUrl = (endpoint: string) => {
    if (!endpoint.startsWith("/")) endpoint = "/" + endpoint
    return PUBLIC_API_URL + endpoint
}

export const getCredentials = async () => {
    const u = await userManager.getUser()
    if (!u) throw new Error("User not logged in")
    const token = u.access_token

    const res = await fetch(getApiUrl("/credentials"), {
        headers: {
            Authorization: `Bearer ${token}`
        }
    });
    if (!res.ok) {
        throw new Error("Failed to get credentials")
    }

    const data = await res.json()
    console.log(data)
}
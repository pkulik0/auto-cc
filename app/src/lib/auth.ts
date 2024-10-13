import { PUBLIC_KEYCLOAK_CLIENT_ID, PUBLIC_KEYCLOAK_URL } from "$env/static/public";
import { UserManager, Log, User, WebStorageStateStore } from "oidc-client-ts";
import { derived, writable } from "svelte/store";

const baseUrl = window.location.origin;
export const userManager = new UserManager({
    authority: PUBLIC_KEYCLOAK_URL,
    client_id: PUBLIC_KEYCLOAK_CLIENT_ID, 
    redirect_uri: `${baseUrl}/auth/callback`,
    post_logout_redirect_uri: `${baseUrl}/`,
    silent_redirect_uri: `${baseUrl}/auth/silent`,
    scope: "openid profile email",
    userStore: new WebStorageStateStore({ store: window.localStorage }),
})

Log.setLogger(console);
Log.setLevel(Log.INFO);

export const userStore = writable<User | null>(null, set => {
    userManager.getUser().then(user => {
        set(user);
    });
    userManager.events.addUserLoaded(user => {
        set(user);
    });
    userManager.events.addUserUnloaded(() => {
        set(null);
    });
})

export const login = async () => {
    await userManager.signinRedirect();
}

export const logout = async () => {
    await userManager.signoutRedirect();
}

const getUserRoles = async (user: User) => {
    const token = user.access_token;
    const parsedToken = JSON.parse(atob(token.split('.')[1]));

    return parsedToken.realm_access?.roles || [];
}

export const isSuperuserStore = derived(userStore, ($user, set) => {
    if (!$user) {
        set(false);
        return;
    }
    getUserRoles($user).then(roles => {
        set(roles.includes("superuser"));
    }).catch((e) => {
        console.error(e);
        set(false);
    });
});
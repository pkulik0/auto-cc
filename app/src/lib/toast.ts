import {writable} from "svelte/store";

export interface Toast {
    id: string
    title: string
    body: string
}

export const toasts = writable<Toast[]>([])
let toastId = 0

const displayToast = async (id: string) => {
    const bootstrap = await import('bootstrap/dist/js/bootstrap.bundle.min.js');
    const toast = new bootstrap.Toast(document.getElementById('toast' + id));
    toast.show();
}

export const sendToast = async (title: string, body: string) => {
    const toast: Toast = {
        "id": (toastId++).toString(),
        "title": title,
        "body": body
    }

    toasts.update(v => {
        v.push(toast)
        return v
    })
    await displayToast(toast.id)

    setInterval(() => {
        toasts.update(v => v.filter(t => t.id !== toast.id))
    }, 10000)
}

export const successOrToast = async (doWork: Function) => {
    try {
        await doWork()
    } catch (e) {
        await sendToast("Error", e.message)
    }
}
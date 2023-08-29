import {alertMessage} from "$lib/alert";

export const successOrAlert = async (doWork: Function) => {
    try {
        await doWork()
    } catch (e) {
        alertMessage.set(e)
        setTimeout( () => {
            let currentMessage = ""
            alertMessage.subscribe(val => currentMessage = val)()
            if(e === currentMessage) alertMessage.set("")
        }, 3000)
    }
}
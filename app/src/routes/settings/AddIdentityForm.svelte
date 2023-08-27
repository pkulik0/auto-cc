<script lang="ts">
    import {addIdentity} from "$lib/youtube/identities";
    import {successOrAlert} from "$lib/error";

    let clientId = ""
    let clientSecret = ""

    $: clientIdTrimmed = clientId.trim()
    $: clientSecretTrimmed = clientSecret.trim()

    let clientIdError = false
    let clientSecretError = false

    const submitForm = async () => {
        clientIdError = !clientIdTrimmed
        clientSecretError = !clientSecretTrimmed
        if(clientIdError || clientSecretError) return

        await successOrAlert(async () => await addIdentity(clientIdTrimmed, clientSecretTrimmed))

        clientId = ""
        clientSecret = ""
    }
</script>

<h5>Add new</h5>
<div class="m-2 mt-4">
    <form on:submit|preventDefault={submitForm}>
        <div class="form-floating mb-3">
            <input bind:value={clientId} type="text" class="form-control" id="clientIdInput" placeholder="Client Id">
            <label for="clientIdInput">Client Id</label>
            {#if clientIdError}<div class="invalid-feedback">Client Id is required</div>{/if}
        </div>

        <div class="form-floating mb-3">
            <input bind:value={clientSecret} type="text" class="form-control" id="clientSecretInput" placeholder="Client Secret">
            <label for="clientSecretInput">Client Secret</label>
            {#if clientSecretError}<div class="invalid-feedback">Client Secret is required</div>{/if}
        </div>

        <button type="submit" class="btn btn-primary w-100">Add</button>
    </form>
</div>

<style>
    ::placeholder {
        opacity: 0.5;
    }
</style>
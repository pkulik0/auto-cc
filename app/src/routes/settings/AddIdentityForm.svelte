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
<div class="m-2">
    <form on:submit|preventDefault={submitForm}>
        <div class="mb-3">
            <label for="clientIdInput" class="form-label">Client Id</label>
            <input on:input={() => clientIdError = false} class:is-invalid={clientIdError} bind:value={clientId} type="text" class="form-control" id="clientIdInput" placeholder="420927583850-arumec536i32pa6v7o0ehqpaqs4ld9e2.apps.googleusercontent.com">
            {#if clientIdError}<div class="invalid-feedback">Client Id is required</div>{/if}
        </div>

        <div class="mb-3">
            <label for="clientSecretInput" class="form-label">Client Secret</label>
            <input on:input={() => clientSecretError = false} class:is-invalid={clientSecretError} bind:value={clientSecret} type="text" class="form-control" id="clientSecretInput" placeholder="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">
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
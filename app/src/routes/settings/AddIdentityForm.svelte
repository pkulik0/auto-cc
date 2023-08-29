<script lang="ts">
    import {addIdentity} from "$lib/youtube/identities";
    import {successOrToast} from "$lib/toast";

    let clientId = ""
    let clientSecret = ""

    let clientIdError = false
    let clientSecretError = false

    const submitForm = async () => {
        const idTrimmed = clientId.trim()
        const secretTrimmed = clientSecret.trim()

        clientIdError = !idTrimmed
        clientSecretError = !secretTrimmed
        if(clientIdError || clientSecretError) return

        await successOrToast(async () => await addIdentity(idTrimmed, secretTrimmed))

        clientId = ""
        clientSecret = ""
    }
</script>

<h5>Add new</h5>
<div class="m-2 mt-4">
    <form on:submit|preventDefault={submitForm}>
        <div class="form-floating mb-3">
            <input on:input={() => clientIdError = false} bind:value={clientId} type="text" class="form-control {clientIdError ? 'is-invalid' : ''}" id="clientIdInput" placeholder="Client Id">
            <label for="clientIdInput">Client Id</label>
            {#if clientIdError}<div class="invalid-feedback">Client Id is required</div>{/if}
        </div>

        <div class="form-floating mb-3">
            <input on:input={() => clientSecretError = false} bind:value={clientSecret} type="text" class="form-control {clientSecretError ? 'is-invalid' : ''}" id="clientSecretInput" placeholder="Client Secret">
            <label for="clientSecretInput">Client Secret</label>
            {#if clientSecretError}<div class="invalid-feedback">Client Secret is required</div>{/if}
        </div>

        <button type="submit" class="btn btn-primary w-100">Add</button>
    </form>
</div>
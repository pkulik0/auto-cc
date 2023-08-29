<script lang="ts">
    import type {VideoMetadata} from "$lib/youtube/metadata";
    import {insertVideoMetadata, metadataSeparator, translateMetadata} from "$lib/youtube/metadata";
    import {onMount} from "svelte";
    import {fade} from "svelte/transition"
    import {sendToast, successOrToast} from "$lib/toast";

    export let videoId: string
    export let metadata: VideoMetadata
    export let hideEditor: () => void

    let title = ""
    let description = ""

    let titleError = false
    let descriptionError = false

    const reset = () => {
        title = metadata.title
        description = metadata.description
    }

    const runTranslation = async () => {
        const editedMetadata: VideoMetadata = {
            "title": title,
            "description": description,
            "language": metadata.language
        }

        await successOrToast(async () => {
            const translatedMetadataArray = await translateMetadata(editedMetadata)
            await insertVideoMetadata(videoId, translatedMetadataArray)

            await sendToast("Success", "The metadata has been updated.")
            hideEditor()
        })
    }

    const submitForm = async () => {
        titleError = !title.trim()
        descriptionError = !description.trim()
        if(titleError || descriptionError) return

        await runTranslation()
    }

    onMount(() => {
        reset()
    })
</script>

<div transition:fade class="row w-75 mx-auto my-2">
    <div class="col-md-8 col-12 order-md-0 order-1">
        <form on:submit|preventDefault={submitForm}>
            <div class="form-floating mb-3">
                <input on:input={() => titleError = false} bind:value={title} type="text" class="form-control {titleError ? 'is-invalid' : ''}" id="titleInput" placeholder="Title">
                <label for="titleInput">Title</label>
                {#if titleError}<div class="invalid-feedback">Title is required</div>{/if}
            </div>

            <div class="form-floating mb-3">
                <textarea on:input={() => descriptionError = false} bind:value={description} class="form-control {descriptionError ? 'is-invalid' : ''}" id="descriptionInput" placeholder="Description" rows="10"></textarea>
                <label for="descriptionInput">Description</label>
                {#if descriptionError}<div class="invalid-feedback">Description is required</div>{/if}
            </div>

            <button type="submit" class="btn btn-primary w-100 m-1">Start translation</button>
            <button type="button" on:click={reset} class="btn btn-outline-secondary w-100 m-1">Revert changes</button>
        </form>
    </div>
    <div class="col-md-4 col-12 mt-2 mt-mb-0 order-md-1 order-0">
        <h5>Example</h5>
        <code>Separators{metadataSeparator} can be used to split{metadataSeparator} text for translation</code>

        <h5 class="mt-3">Resulting translation input</h5>
        <p>1.<code class="ms-2">Separators</code></p>
        <p>2.<code class="ms-2"> can be used to split</code></p>
        <p>3.<code class="ms-2"> text for translation</code></p>
    </div>
</div>

<style>
    #descriptionInput {
        height: 15em;
    }
</style>
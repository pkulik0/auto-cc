<script lang="ts">
    import type {Video} from "$lib/youtube/video";
    import {filteredTargetLanguages, selectedLanguage} from "$lib/languages/data";
    import {translateVideoCC} from "$lib/youtube/cc";
    import MetadataEditor from "./MetadataEditor.svelte";
    import type {VideoMetadata} from "$lib/youtube/metadata";
    import {getMetadata} from "$lib/youtube/metadata";
    import {sendToast, successOrToast} from "$lib/toast";

    export let video: Video
    let videoTime = new Date(video.publishedAt).toLocaleString()

    let showMetadataEditor = false
    let metadata: VideoMetadata

    $: metadataButtonLabel = (showMetadataEditor ? "Hide" : "Show") + " metadata translation"

    const runCaptionsTranslation = async () => {
        if(!$selectedLanguage) {
            alert("No source language selected.")
            return
        }

        await successOrToast(async () => {
            await translateVideoCC(video, $selectedLanguage!.language, filteredTargetLanguages)
            await sendToast("Success", "The closed captions have been uploaded to YouTube.")
        })
    }

    const toggleMetadataEditor = async () => {
        await successOrToast(async () => {
            metadata = await getMetadata(video.id)
            if(!metadata.language) throw new Error("Unknown default language. Set in on YouTube first.")

            showMetadataEditor = !showMetadataEditor
        })
    }
</script>

<div class="row align-items-center text-center justify-content-center">
    <div class="col-md-3 col-12"><img alt="" class="w-75 p-3" src={video.thumbnailUrl}></div>
    <div class="col-md-4 mb-md-0 col-12 mb-2 lead">{video.title}</div>
    <div class="col-md-2 mb-md-0 col-12 mb-2">{videoTime}</div>
    <div class="col-md-3 mb-md-0 col-12 mb-2">
        <button on:click={runCaptionsTranslation} class="btn btn-primary w-100 m-1">Translate captions</button>
        <button on:click={toggleMetadataEditor} class="btn btn-outline-primary w-100 m-1">{metadataButtonLabel}</button>
    </div>
</div>

{#if showMetadataEditor}
    <MetadataEditor videoId={video.id} {metadata} hideEditor={() => showMetadataEditor = false}/>
{/if}

<hr class="hr" />
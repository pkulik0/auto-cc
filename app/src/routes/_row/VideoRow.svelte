<script lang="ts">
    import type {Video} from "$lib/youtube/video";
    import {filteredTargetLanguages, selectedLanguage, targetLanguages} from "$lib/languages/data";
    import {translateVideoCC} from "$lib/youtube/cc";
    import {successOrAlert} from "$lib/error";
    import MetadataEditor from "./MetadataEditor.svelte";
    import type {VideoMetadata} from "$lib/youtube/metadata";
    import {getMetadata} from "$lib/youtube/metadata";

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

        await successOrAlert(async () => await translateVideoCC(video, $selectedLanguage!.language, filteredTargetLanguages))
    }

    const toggleMetadataEditor = async () => {
        await successOrAlert(async () => {
            metadata = await getMetadata(video.id)
            if(!metadata.language) throw new Error("Unknown default language. Set in on YouTube first.")

            showMetadataEditor = !showMetadataEditor
        })
    }
</script>

<div class="row align-items-center text-center justify-content-center">
    <div class="col-3"><img alt="" class="w-75 p-3" src={video.thumbnailUrl}></div>
    <div class="col-4 lead">{video.title}</div>
    <div class="col-2">{videoTime}</div>
    <div class="col-3">
        <button on:click={runCaptionsTranslation} class="btn btn-primary w-100 m-1">Translate captions</button>
        <button on:click={toggleMetadataEditor} class="btn btn-outline-primary w-100 m-1">{metadataButtonLabel}</button>
    </div>
</div>

{#if showMetadataEditor}
    <MetadataEditor videoId={video.id} {metadata} />
{/if}
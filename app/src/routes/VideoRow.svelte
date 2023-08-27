<script lang="ts">
    import type {Video} from "$lib/youtube/video";
    import {filteredTargetLanguages, selectedLanguage, targetLanguages} from "$lib/languages/data";
    import {translateVideoCC} from "$lib/youtube/cc";
    import {successOrAlert} from "$lib/error";

    export let video: Video
    let videoTime = new Date(video.publishedAt).toLocaleString()

    const translateCaptions = async () => {
        if(!$selectedLanguage) {
            alert("No source language selected.")
            return
        }

        await successOrAlert(async () => await translateVideoCC(video, $selectedLanguage, filteredTargetLanguages))
    }

    const translateMetadata = async () => {

    }
</script>

<tr class="align-middle">
    <td class="w-25"><img alt="" class="w-75 p-3" src={video.thumbnailUrl}></td>
    <td class="lead">{video.title}</td>
    <td>{videoTime}</td>
    <td class="w-25">
        <button on:click={translateMetadata} class="btn btn-outline-primary w-100 m-1">Translate metadata</button>
        <button on:click={translateCaptions} class="btn btn-primary w-100 m-1">Translate captions</button>
    </td>
</tr>
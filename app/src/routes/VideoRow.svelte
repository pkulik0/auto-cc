<script lang="ts">
    import type {CCEntry, Video} from "$lib/youtube/video";
    import {selectedLanguage} from "$lib/languages/stores";
    import {translateVideoCC} from "$lib/youtube/cc";

    export let video: Video
    let videoTime = new Date(video.publishedAt*1000).toLocaleString()

    const translate = async () => {
        if(!$selectedLanguage) {
            alert("No source language selected.")
            return
        }
        await translateVideoCC($selectedLanguage, video)
    }
</script>

<tr>
    <td><img alt="" width="300" src={video.thumbnailUrl}></td>
    <td>{video.id}</td>
    <td>{video.title}</td>
    <td>{videoTime}</td>
    <td>
        <button on:click={translate} class="btn btn-primary w-100">
            Translate
        </button>
    </td>
</tr>
<script lang="ts">
    import type {Video} from "$lib/youtube/video";
    import {selectedLanguage} from "$lib/languages/data";
    import {translateVideoCC} from "$lib/youtube/cc";

    export let video: Video
    let videoTime = new Date(video.publishedAt).toLocaleString()

    const translate = async () => {
        if(!$selectedLanguage) {
            alert("No source language selected.")
            return
        }
        await translateVideoCC($selectedLanguage, video)
    }
</script>

<tr class="align-middle">
    <td><img alt="" class="w-75 p-3" src={video.thumbnailUrl}></td>
    <td class="lead">{video.title}</td>
    <td>{videoTime}</td>
    <td>
        <button on:click={translate} class="btn btn-primary w-100">
            Translate
        </button>
    </td>
</tr>
<script lang="ts">
    import {getVideos, videosNextPageToken, videos} from "$lib/youtube/video";
    import type {Video} from "$lib/youtube/video";
    import VideoRow from "./_row/VideoRow.svelte";
    import {onMount} from "svelte";
    import {successOrAlert} from "$lib/error";

    const loadVideos = async () => {
        await successOrAlert(async () => videos.set(await getVideos()))
    }

    onMount(async () => {
        await loadVideos()
    })

    const loadMoreVideos = async () => {
        const moreVideos = await getVideos(true, true)
        videos.update((v: Video[]) => [...v, ...moreVideos])
    }
</script>

{#each $videos as video}
    <div class="mb-3">
        <VideoRow {video}/>
    </div>
{/each}

{#if $videosNextPageToken}
    <button on:click={loadMoreVideos} class="btn btn-outline-primary w-25 mx-auto d-block mb-4">Load more</button>
{/if}
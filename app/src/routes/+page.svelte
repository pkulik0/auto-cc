<script lang="ts">
    import {getVideos, videosNextPageToken, videos} from "$lib/youtube/video";
    import type {Video} from "$lib/youtube/video";
    import VideoRow from "./_row/VideoRow.svelte";
    import {onMount} from "svelte";
    import {successOrToast} from "$lib/toast";

    const loadVideos = async () => {
        await successOrToast(async () => videos.set(await getVideos()))
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
    <button on:click={loadMoreVideos} class="btn btn-outline-primary w-75 mx-auto d-block mb-4">Load more</button>
{/if}
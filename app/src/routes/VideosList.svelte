<script lang="ts">
    import {getVideos, nextPageToken, videos} from "$lib/youtube/video";
    import type {Video} from "$lib/youtube/video";
    import VideoRow from "./VideoRow.svelte";
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

<div class="table-responsive mt-4">
    <table class="table text-center">
        <tbody>
        {#each $videos as video}
            <VideoRow {video}/>
        {/each}
        </tbody>
    </table>
    {#if $nextPageToken}
        <button on:click={loadMoreVideos} class="btn btn-outline-primary w-25 mx-auto d-block mb-4">Load more</button>
    {/if}
</div>
<script lang="ts">
    import {getVideos, nextPageToken, videos} from "$lib/youtube/video";
    import type {Video} from "$lib/youtube/video";
    import VideoRow from "./VideoRow.svelte";
    import {onMount} from "svelte";

    onMount(async () => {
        videos.set(await getVideos())
    })

    const loadMoreVideos = async () => {
        const moreVideos = await getVideos(true, true)
        videos.update((v: Video[]) => [...v, ...moreVideos])
    }
</script>

<div class="table-responsive">
    <table class="table table-striped text-center">
        <thead class="fs-5">
        <tr>
            <th>Thumbnail</th>
            <th>ID</th>
            <th>Title</th>
            <th>Published at</th>
            <th>Actions</th>
        </tr>
        </thead>
        <tbody>
            {#each $videos as video}
                <VideoRow {video}/>
            {/each}
        </tbody>
    </table>
    {#if $nextPageToken}
        <button on:click={loadMoreVideos} class="btn btn-primary mx-auto d-block mb-4">Load more</button>
    {/if}
</div>
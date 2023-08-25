<script lang="ts">
    import {getVideos, videos} from "$lib/youtube/video";
    import VideoRow from "./VideoRow.svelte";
    import {onMount} from "svelte";

    onMount(async () => {
        videos.set(await getVideos())
    })
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
</div>
<script lang="ts">
    import type {Video} from "$lib/youtube";
    import {onMount} from "svelte";
    import {getVideos} from "$lib/youtube";

    let videos: Video[] = []

    onMount(async () => {
        videos = await getVideos()
        console.log(videos)
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
        {#each videos as video}
            <tr>
                <td><img alt="" src={video.thumbnailUrl} width="300"></td>
                <td>{video.id}</td>
                <td>{video.title}</td>
                <td>{new Date(video.publishedAt*1000).toLocaleString()}</td>
                <td>
                    <button class="btn btn-primary w-100" on:click={() => console.log(video)}>
                        Translate
                    </button>
                </td>
            </tr>
        {/each}
        </tbody>
    </table>
</div>
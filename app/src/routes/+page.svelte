<script lang="ts">
    import {getVideos} from "$lib/youtube/api";
    import VideoRow from "./VideoRow.svelte";

    let videosPromise = getVideos()
</script>

{#await videosPromise}
    <div class="d-flex justify-content-center">
        <div class="spinner-border m-5" role="status"/>
    </div>
{:then videos}
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
                    <VideoRow {video}/>
                {/each}
            </tbody>
        </table>
    </div>
{:catch error}
    {error}
{/await}
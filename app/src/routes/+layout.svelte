<script>
    import LanguageDropdown from "./LanguageDropdown.svelte";
    import {getVideos, videos} from "$lib/youtube/video";

    const refreshVideos = async () => {
        if(!confirm("Are you sure?")) return
        videos.set(await getVideos(true))
    }
</script>

<nav class="navbar navbar-expand-lg bg-primary mb-2">
    <div class="container-fluid">
        <div class="navbar-brand">AutoCC</div>

        <div class="navbar-nav">
            <button class="btn btn-warning m-1" on:click={refreshVideos}>Refresh videos</button>
            <LanguageDropdown/>
        </div>
    </div>
</nav>

<div class="container-fluid">
    <slot/>
</div>
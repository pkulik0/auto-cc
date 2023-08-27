<script>
    import LanguageDropdown from "./LanguageDropdown.svelte";
    import {getVideos, videos} from "$lib/youtube/video";
    import {page} from "$app/stores";
    import {routes} from "$lib/routes";

    const refreshVideos = async () => {
        if(!confirm("Are you sure?")) return
        videos.set(await getVideos(true))
    }
</script>


<nav class="navbar navbar-expand-lg bg-primary mb-2">
    <div class="container">
        <div class="navbar-brand">AutoCC</div>

        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
            <span class="navbar-toggler-icon"></span>
        </button>

        <div class="collapse navbar-collapse" id="navbarNav">
            <ul class="navbar-nav">
                {#each routes as route}
                    <li class="nav-item">
                        <a class:active={$page.route.id === route.destination} class="nav-link" href={route.destination}>{route.label}</a>
                    </li>
                {/each}
            </ul>
            <div class="d-flex ms-auto">
                <button class="btn btn-secondary m-1" on:click={refreshVideos}>Refresh videos</button>
                <LanguageDropdown/>
            </div>
        </div>
    </div>
</nav>
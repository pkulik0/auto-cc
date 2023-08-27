<script lang="ts">
    import "../app.scss"

    import NavBar from "./_nav/NavBar.svelte";
    import {onMount} from "svelte";
    import Alert from "./Alert.svelte";

    onMount(() => {
        import("bootstrap/dist/js/bootstrap.bundle.min.js");
    });
</script>

<NavBar/>

<div class="container">
    <Alert/>
    <slot/>
</div>
<script lang="ts">
    import "../app.scss"

    import NavBar from "./_nav/NavBar.svelte";
    import {onMount} from "svelte";
    import ToastContainer from "./ToastContainer.svelte";

    onMount(() => {
        import("bootstrap/dist/js/bootstrap.bundle.min.js");
    });
</script>

<ToastContainer/>
<NavBar/>

<div class="container">
    <slot/>
</div>
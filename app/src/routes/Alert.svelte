<script lang="ts">
    import {fade} from "svelte/transition"
    import {alertMessage} from "$lib/alert";
</script>

{#if $alertMessage}
    <div transition:fade class="alert alert-warning my-4" role="alert">
        {$alertMessage}
    </div>
{/if}

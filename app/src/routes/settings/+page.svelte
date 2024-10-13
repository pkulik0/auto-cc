<script lang="ts">
	import { isSuperuserStore, userStore } from "$lib/auth";
	import { onMount } from "svelte";

</script>

{#if $isSuperuserStore}
abc
{:else}
def
{/if}
<script lang="ts">
	import { getVideos } from '$lib/api';
	import { isSuperuserStore } from '$lib/auth';
	import { onMount } from 'svelte';

	onMount(async () => {
		try {
			const res = await getVideos();
			console.log(res);
		} catch (error) {
			console.error(error);
		}
	});
</script>

{#if $isSuperuserStore}
	abc
{:else}
	def
{/if}

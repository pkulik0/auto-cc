<script lang="ts">
	import { getMetadata, getVideos } from '$lib/api';
	import { isSuperuserStore } from '$lib/auth';
	import { onMount } from 'svelte';

	onMount(async () => {
		try {
			const res = await getVideos();
			console.log(res);

			const res2 = await getMetadata(res.videos[0].id);
			console.log(res2);
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

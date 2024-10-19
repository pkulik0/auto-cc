<script lang="ts">
	import { getCC, getMetadata, getVideos } from '$lib/api';
	import { isSuperuserStore } from '$lib/auth';
	import { onMount } from 'svelte';

	onMount(async () => {
		try {
			const res = await getVideos();
			console.log(res);

			const id = res.videos[4].id;

			const res2 = await getMetadata(id);
			console.log(res2);

			const res3 = await getCC(id);
			console.log(res3);
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

<script lang="ts">
	import { downloadCC, getCC, getLanguages, getMetadata, getVideos, translate } from '$lib/api';
	import { isSuperuserStore } from '$lib/auth';
	import { onMount } from 'svelte';

	onMount(async () => {
		try {
			// const res = await getVideos();
			// console.log(res);

			// const id = res.videos[4].id;

			// const res2 = await getMetadata(id);
			// console.log(res2);

			// const res3 = await getCC(id);
			// console.log(res3);

			// const res4 = await downloadCC(res3[0].id);
			// console.log(res4);

			const res0 = await getLanguages();
			console.log(res0);

			const res = await translate(["good morning, how are you?", "very well, thanks for asking"], "en", "pl")
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

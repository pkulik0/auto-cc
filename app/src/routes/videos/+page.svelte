<script lang="ts">
	import { Card, Button, Spinner } from 'flowbite-svelte';
	import { getLanguages, getVideos, process } from '$lib/api';
	import type { Video } from '$lib/pb/youtube';
	import { onMount } from 'svelte';
	import { _ } from 'svelte-i18n';
	import { CaptionOutline, ListOutline } from 'flowbite-svelte-icons';
	import { fade } from 'svelte/transition';

	let videos: Video[] | null = null;
	let videosNextPageToken: string = '';

	let sentiel: Element;
	let isLoading = false;

	const fetch = async () => {
		isLoading = true;
		try {
			const videosResp = await getVideos(videosNextPageToken);
			videosNextPageToken = videosResp.nextPageToken;
			if (videos) {
				videos = [...videos, ...videosResp.videos];
			} else {
				videos = videosResp.videos;
			}
		} catch (error) {
			console.error(error);
		}
		isLoading = false;
	};

	onMount(() => {
		fetch();

		const observer = new IntersectionObserver(
			(entries) => {
				if (!videosNextPageToken || !videos || isLoading) {
					return;
				}

				if (entries[0].isIntersecting) fetch();
			},
			{
				root: null,
				rootMargin: '0px',
				threshold: 1.0
			}
		);
		observer.observe(sentiel);

		return () => {
			observer.disconnect();
		};
	});
</script>

{#if videos}
	{#if videos.length === 0}
		<div class="flex flex-row space-x-4">
			<span class="mx-auto mt-8 text-center text-gray-700 dark:text-gray-400">
				<ListOutline class="mx-auto h-12 w-12 text-gray-400 dark:text-gray-500" />
				{$_('videos.no_videos')}
			</span>
		</div>
	{/if}
	<div
		class="mx-auto grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4"
		transition:fade
	>
		{#each videos as video}
			<div transition:fade>
				<Card class="flex h-full flex-col">
					<img
						src={video.thumbnailUrl}
						alt={video.title}
						class="mb-4 aspect-video rounded-md object-cover"
					/>
					<div class="flex-grow">
						<h5
							class="mb-2 truncate text-lg font-bold tracking-tight text-gray-900 dark:text-white"
						>
							{video.title}
						</h5>

						<p class="line-clamp-3 text-sm leading-tight text-gray-700 dark:text-gray-400">
							{#if video.description}
								{video.description}
							{:else}
								{$_('videos.no_description')}
							{/if}
						</p>
					</div>

					<div class="mt-8 space-y-2">
						<Button on:click={() => process(video.id)} class="w-full" size="xs">
							<CaptionOutline class="me-1 w-6" />
							{$_('videos.closed_captions')}
						</Button>
						<!-- <Button class="w-full" outline size="xs">{$_('videos.metadata')}</Button> -->
					</div>
				</Card>
			</div>
		{/each}
	</div>
{/if}

{#if !videos || isLoading}
	<div class="mx-auto mt-8 flex justify-center">
		<Spinner />
	</div>
{/if}

<div bind:this={sentiel}></div>

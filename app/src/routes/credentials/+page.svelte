<script lang="ts">
	import { onMount } from 'svelte';
	import { isSuperuserStore } from '$lib/auth';
	import { Button, Input, Select, Hr, Skeleton } from 'flowbite-svelte';
	import { PlusOutline } from 'flowbite-svelte-icons';
	import {
		addCredentialsDeepL,
		addCredentialsGoogle,
		getCredentials,
		type Credentials
	} from '$lib/api';
	import GoogleTable from './GoogleTable.svelte';
	import DeepLTable from './DeepLTable.svelte';
	import { _ } from 'svelte-i18n';
	import { fade } from 'svelte/transition';

	let credentials: Credentials | null = null;

	onMount(async () => {
		try {
			credentials = await getCredentials();
		} catch (error) {
			console.error(error);
		}
	});

	const typeItems = [
		{ value: 'deepl', name: 'DeepL' },
		{ value: 'google', name: 'Google' }
	];
	let selectedType = typeItems[0].value;

	let client_id = '';
	let secret = '';
	$: secretPlaceholder =
		selectedType === 'google' ? $_('credentials.client_secret') : $_('credentials.api_key');

	$: isAddEnabled = selectedType === 'google' ? client_id.trim() && secret.trim() : secret.trim();

	const addCredential = async () => {
		console.log(selectedType, client_id, secret);

		try {
			switch (selectedType) {
				case 'google':
					const googleCred = await addCredentialsGoogle(client_id, secret);
					if (!credentials) return;
					credentials = { ...credentials, google: [...credentials.google, googleCred] };
					break;
				case 'deepl':
					const deeplCred = await addCredentialsDeepL(secret);
					if (!credentials) return;
					credentials = { ...credentials, deepl: [...credentials.deepl, deeplCred] };
					break;
			}
			client_id = '';
			secret = '';
		} catch (error) {
			console.error(error);
		}
	};
</script>

<div class="space-y-8">
	{#if $isSuperuserStore}
		<h2 class="mb-4 text-lg font-semibold text-gray-800 dark:text-gray-200">Add new credential</h2>
		<div class="flex space-x-4">
			<Select placeholder={$_('credentials.service')} items={typeItems} bind:value={selectedType} />
			{#if selectedType === 'google'}
				<Input placeholder={$_('credentials.client_id')} bind:value={client_id} />
			{/if}
			<Input placeholder={secretPlaceholder} bind:value={secret} />
			<Button size="xs" on:click={addCredential} disabled={!isAddEnabled}>
				<PlusOutline class="h-6 w-6" />
			</Button>
		</div>
		<Hr />
	{/if}

	{#if credentials}
		<div class="space-y-10" transition:fade>
			<GoogleTable credentials={credentials.google} />
			<DeepLTable credentials={credentials.deepl} />
		</div>
	{:else}
		<div class="space-y-10">
			<Skeleton size="80" />
			<Skeleton size="80" />
		</div>
	{/if}
</div>

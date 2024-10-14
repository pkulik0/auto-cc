<script lang="ts">
	import type { CredentialsGoogle } from '$lib/pb/autocc';
	import {
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell,
		Button,
		Progressbar
	} from 'flowbite-svelte';
	import { TrashBinOutline } from 'flowbite-svelte-icons';
	import { _ } from 'svelte-i18n';
	import {
		getSessionGoogleURL,
		getUserSessionsGoogle,
		removeCredentials,
		removeSessionGoogle
	} from '$lib/api';
	import { isSuperuserStore } from '$lib/auth';
	import { onMount } from 'svelte';

	const quota = 10000; // Youtube API quota

	export let credentials: CredentialsGoogle[];
	let sessions: number[] = [];

	onMount(async () => {
		try {
			sessions = await getUserSessionsGoogle();
		} catch (error) {
			console.error(error);
		}
	});

	const remove = async (id: number) => {
		try {
			await removeCredentials('google', id);
			credentials = credentials.filter((c) => c.id !== id);
		} catch (error) {
			console.error(error);
		}
	};

	const authenticate = async (id: number) => {
		const url = await getSessionGoogleURL(id);
		window.location.href = url;
	};

	const revoke = async (id: number) => {
		try {
			await removeSessionGoogle(id);
			sessions = sessions.filter((s) => s !== id);
		} catch (error) {
			console.error(error);
		}
	};
</script>

<h2 class="mb-4 text-xl font-semibold text-gray-800 dark:text-gray-200">Google</h2>

<Table striped={true}>
	<TableHead>
		<TableHeadCell>{$_('credentials.client_id')}</TableHeadCell>
		<TableHeadCell>{$_('credentials.client_secret')}</TableHeadCell>
		<TableHeadCell>{$_('credentials.usage')}</TableHeadCell>
		<TableHeadCell>{$_('credentials.actions')}</TableHeadCell>
	</TableHead>
	<TableBody tableBodyClass="divide-y">
		{#if credentials.length === 0}
			<TableBodyRow>
				<TableBodyCell colspan={4} class="text-center">
					{$_('credentials.no_credentials')}
				</TableBodyCell>
			</TableBodyRow>
		{/if}
		{#each credentials as credential}
			<TableBodyRow>
				<TableBodyCell>{credential.clientId.substring(0, 30)}...</TableBodyCell>
				<TableBodyCell>{credential.clientSecret}</TableBodyCell>
				<TableBodyCell>
					<Progressbar class="w-60" progress={credential.usage * 100 / quota} />
				</TableBodyCell>
				<TableBodyCell class="flex items-center space-x-2">
					{#if sessions.includes(credential.id)}
						<Button size="xs" outline on:click={() => revoke(credential.id)}>
							{$_('credentials.revoke')}
						</Button>
					{:else}
						<Button size="xs" outline color="green" on:click={() => authenticate(credential.id)}>
							{$_('credentials.authenticate')}
						</Button>
					{/if}

					{#if $isSuperuserStore}
						<Button size="xs" outline color="red" on:click={() => remove(credential.id)}>
							<TrashBinOutline class="h-4 w-4" />
						</Button>
					{/if}
				</TableBodyCell>
			</TableBodyRow>
		{/each}
	</TableBody>
</Table>

<script lang="ts">
	import type { CredentialsGoogle } from '$lib/pb/autocc';
	import {
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell,
		Button
	} from 'flowbite-svelte';
	import { TrashBinOutline } from 'flowbite-svelte-icons';
	import { _ } from 'svelte-i18n';
	import { removeCredentials } from '$lib/api';
	import { isSuperuserStore } from '$lib/auth';

	export let credentials: CredentialsGoogle[];

	const remove = async (id: number) => {
		try {
			await removeCredentials("google", id)
			credentials = credentials.filter((c) => c.id !== id);
		} catch (error) {
			console.error(error);
		}
	};
</script>

<h2 class="mb-4 text-xl font-semibold text-gray-800 dark:text-gray-200">
	Google
</h2>

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
				<TableBodyCell class="w-1/4">{credential.clientId}</TableBodyCell>
				<TableBodyCell class="w-1/4">{credential.clientSecret}</TableBodyCell>
				<TableBodyCell class="w-1/4">{credential.usage}</TableBodyCell>
				<TableBodyCell class="w-1/4">
					{#if $isSuperuserStore}
						<Button size="xs" outline color="red" on:click={() => remove(credential.id)}>
							<TrashBinOutline class="w-4 h-4" />
						</Button>
					{/if}
				</TableBodyCell>
			</TableBodyRow>
		{/each}
	</TableBody>
</Table>

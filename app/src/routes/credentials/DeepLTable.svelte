<script lang="ts">
	import { removeCredentials } from '$lib/api';
	import { isSuperuserStore } from '$lib/auth';
	import type { CredentialsDeepL } from '$lib/pb/credentials';
	import { QuotaDeepL } from '$lib/quota';
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

	export let credentials: CredentialsDeepL[];

	const remove = async (id: number) => {
		try {
			await removeCredentials('deepl', id);
			credentials = credentials.filter((c) => c.id !== id);
		} catch (error) {
			console.error(error);
		}
	};
</script>

<h2 class="mb-4 text-xl font-semibold text-gray-800 dark:text-gray-200">DeepL</h2>

<Table striped={true}>
	<TableHead>
		<TableHeadCell colspan={2}>{$_('credentials.api_key')}</TableHeadCell>
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
				<TableBodyCell class="w-1/2" colspan={2}>{credential.key}</TableBodyCell>
				<TableBodyCell>
					<Progressbar class="w-60" progress={(credential.usage * 100) / QuotaDeepL} />
				</TableBodyCell>
				<TableBodyCell class="w-1/4">
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

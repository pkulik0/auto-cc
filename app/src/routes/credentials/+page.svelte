<script>
	import { onMount } from "svelte";
	import CredentialsTable from "./CredentialsTable.svelte";
	import { getCredentials } from "$lib/api";

    onMount(async () => {
        try {
            await getCredentials();
            console.log("Credentials fetched");
        } catch (error) {
            console.error(error);
        }
    })

</script>

<h1 class="mb-4 text-xl font-semibold text-gray-800 dark:text-gray-200">Credentials</h1>

<div class="space-y-8 m-10">
    <CredentialsTable service="Google"/>

    <CredentialsTable service="DeepL"/>
</div>
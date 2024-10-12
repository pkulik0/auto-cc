<script>
	import CredentialsTable from "./CredentialsTable.svelte";

</script>

<h1 class="mb-4 text-xl font-semibold text-gray-800 dark:text-gray-200">Credentials</h1>

<div class="space-y-8 m-10">
    <CredentialsTable service="Google"/>

    <CredentialsTable service="DeepL"/>
</div>
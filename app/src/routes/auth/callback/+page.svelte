<script lang="ts">
	import { goto } from "$app/navigation";
	import { userManager } from "$lib/auth";
	import { Spinner } from "flowbite-svelte";
	import { onMount } from "svelte";

    onMount(async () => {
        try {
            await userManager.signinRedirectCallback()
        } catch (error) {
            console.error('Error processing login:', error)
        }
        goto('/')
    })
</script>

<div class="flex items-center justify-center">
    <Spinner size={16} />
</div>

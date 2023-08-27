<script lang="ts">
    import type {Identity} from "$lib/youtube/identities";
    import {onDestroy, onMount} from "svelte";
    import {getIdentities} from "$lib/youtube/identities";
    import {successOrAlert} from "$lib/error";

    let identities: Identity[] = []
    let intervalId: number

    const loadIdentities = async () => {
        identities = await getIdentities()
        intervalId = setInterval(async () => {
            identities = await getIdentities()
        }, 5000)
    }

    onMount(async () => {
        await successOrAlert(async () => await loadIdentities())
    })

    onDestroy(() => {
        clearInterval(intervalId)
    })
</script>

<h5>Identities</h5>
<div class="m-2">
    {#each identities as identity}
        <span class:fw-bolder={identity.isSelected}>{identity.isSelected ? "*" : ""}{identity.identityHash.slice(0, 16)}</span>
        <div class="progress mb-4 mt-1" role="progressbar" style="height: 20px;">
            <div class="progress-bar" style="width: {identity.usedQuota / 100}%">{identity.usedQuota}</div>
        </div>
    {/each}
</div>
<script lang="ts">
    import {onDestroy, onMount} from "svelte";
    import type {IdentityInfo} from "$lib/youtube/identities";
    import {getIdentityInfos} from "$lib/youtube/identities";
    import {successOrToast} from "$lib/toast";

    let identityInfos: IdentityInfo[] = []
    let intervalId: number

    const loadIdentities = async () => {
        identityInfos = await getIdentityInfos()
        intervalId = setInterval(async () => {
            identityInfos = await getIdentityInfos()
        }, 5000)
    }

    onMount(async () => {
        await successOrToast(async () => await loadIdentities())
    })

    onDestroy(() => {
        clearInterval(intervalId)
    })
</script>

<h5>Identities</h5>
<div class="m-2">
    {#each identityInfos as identity}
        <div class="my-2">
            <span class:fw-bolder={identity.isSelected}>
                {identity.isSelected ? "*" : ""}{identity.hash.slice(0, 16)}
            </span>
            {#if identity.authUrl}
                <a href={identity.authUrl} target="_blank"><button class="btn btn-outline-primary w-100">Click to authenticate</button></a>
                <br/>
            {:else}
                <div class="progress mb-4 mt-1" role="progressbar" style="height: 35px;">
                    <div class="progress-bar overflow-visible" style="width: {identity.usedQuota / 100}%">
                        <span class="mx-2">{identity.usedQuota}</span>
                    </div>
                </div>
            {/if}
        </div>
    {/each}
</div>
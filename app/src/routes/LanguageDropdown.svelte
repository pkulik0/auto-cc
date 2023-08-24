<script lang="ts">
    import {onMount} from "svelte";
    import {getLanguages} from "$lib/languages";
    import type {Language} from "$lib/languages";

    let sourceLanguages: Language[] = []
    let selectedLanguage: Language|null = null
    $: sourceName = selectedLanguage !== null ? selectedLanguage.name : "..."

    const saveLanguage = (language: Language) => {
        selectedLanguage = language
    }

    onMount(async () => {
        sourceLanguages = (await getLanguages()).source
    })
</script>

<div class="dropdown">
    <button class="btn btn-outline-warning dropdown-toggle" type="button" id="dropdownMenuButton" data-bs-toggle="dropdown" aria-expanded="false">
        Translate from {sourceName}
    </button>
    <ul class="dropdown-menu" aria-labelledby="dropdownMenuButton">
        {#each sourceLanguages as language}
            <li><a class="dropdown-item" on:click={() => saveLanguage(language)} href="/">{language.name}</a></li>
        {/each}
    </ul>
</div>
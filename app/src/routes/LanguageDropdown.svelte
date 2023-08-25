<script lang="ts">
    import {onMount} from "svelte";
    import {getLanguages} from "$lib/languages/api";
    import type {Language} from "$lib/languages/api";
    import {selectedLanguage} from "$lib/languages/data";
    import {setTargetLanguages} from "$lib/languages/data";

    onMount(async () => {
        const response = await getLanguages()
        sourceLanguages = response.source
        setTargetLanguages(response.target)
    })

    const saveLanguage = (language: Language) => {
        selectedLanguage.set(language)
    }

    let sourceLanguages: Language[] = []
    $: sourceLanguageName = $selectedLanguage !== null ? $selectedLanguage.name : "..."

</script>

<div class="dropdown">
    <button class="btn btn-outline-warning dropdown-toggle" type="button" id="dropdownMenuButton" data-bs-toggle="dropdown" aria-expanded="false">
        Translate from {sourceLanguageName}
    </button>
    <ul class="dropdown-menu" aria-labelledby="dropdownMenuButton">
        {#each sourceLanguages as language}
            <li><a class="dropdown-item" on:click={() => saveLanguage(language)} href="/">{language.name}</a></li>
        {/each}
    </ul>
</div>
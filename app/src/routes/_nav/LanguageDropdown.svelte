<script lang="ts">
    import {onMount} from "svelte";
    import {getLanguages} from "$lib/languages/api";
    import type {Language} from "$lib/languages/api";
    import {selectedLanguage} from "$lib/languages/data";
    import {setTargetLanguages} from "$lib/languages/data";
    import {successOrToast} from "$lib/toast";

    const initLanguages = async () => {
        const response = await getLanguages()
        sourceLanguages = response.source
        setTargetLanguages(response.target)
    }

    onMount(async () => {
        await successOrToast(async () => await initLanguages())
    })

    const saveLanguage = (language: Language) => {
        selectedLanguage.set(language)
    }

    let sourceLanguages: Language[] = []
    $: sourceLanguageName = $selectedLanguage !== null ? $selectedLanguage.name : "..."

</script>

<div class="dropdown m-1">
    <button class="btn btn-outline-secondary dropdown-toggle" type="button" id="dropdownMenuButton" data-bs-toggle="dropdown" aria-expanded="false">
        Translate from {sourceLanguageName}
    </button>
    <ul class="dropdown-menu" aria-labelledby="dropdownMenuButton">
        {#each sourceLanguages as language}
            <li><button class="dropdown-item" on:click={() => saveLanguage(language)}>{language.name}</button></li>
        {/each}
    </ul>
</div>
<script lang="ts">
    import type {CCEntry, Video} from "$lib/youtube/api";
    import {selectedLanguage} from "$lib/languages/stores";
    import {downloadCC, getCCList} from "$lib/youtube/api";
    import {Srt, translateSrt} from "$lib/srt";

    export let video: Video
    let videoTime = new Date(video.publishedAt*1000).toLocaleString()

    const translate = async () => {
        if(!$selectedLanguage) {
            throw new Error("No source language selected.")
        }
        const language: string = $selectedLanguage.language.toLowerCase()

        const ccList: CCEntry[] = await getCCList(video.id)
        if(!ccList) {
            throw new Error("No CCs found.")
        }

        const ccEntry: CCEntry = ccList.find(cc => cc.language === language)
        if(!ccEntry) {
            throw new Error(`${$selectedLanguage.language} CC not found!`)
        }

        const srt: Srt = new Srt(await downloadCC(ccEntry.id))
        const translatedSrts: Srt[] = await translateSrt(srt)
        console.log(translatedSrts)
    }
</script>

<tr>
    <td><img alt="" width="300" src={video.thumbnailUrl}></td>
    <td>{video.id}</td>
    <td>{video.title}</td>
    <td>{videoTime}</td>
    <td>
        <button on:click={translate} class="btn btn-primary w-100">
            Translate
        </button>
    </td>
</tr>
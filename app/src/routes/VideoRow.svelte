<script lang="ts">
    import type {CCEntry, Video} from "$lib/youtube/api";
    import {selectedLanguage} from "$lib/languages/stores";
    import {downloadCC, getCCList} from "$lib/youtube/api";
    import {Srt} from "$lib/srt";

    export let video: Video
    let videoTime = new Date(video.publishedAt*1000).toLocaleString()

    const translate = async () => {
        if(!$selectedLanguage) {
            alert("Choose a source language!")
            return
        }
        const language: string = $selectedLanguage.language.toLowerCase()

        const ccList: CCEntry[] = await getCCList(video.id)
        if(!ccList) {
            alert(`No CCs found!`)
            return
        }

        const ccEntry: CCEntry = ccList.find(cc => cc.language === language)
        if(!ccEntry) {
            alert(`${$selectedLanguage.language} CC not found!`)
        }

        const srt: Srt = new Srt(await downloadCC(ccEntry.id))
        console.log(srt)
        console.log(srt.toString())
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
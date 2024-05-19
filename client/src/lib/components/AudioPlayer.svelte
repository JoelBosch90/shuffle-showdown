<script lang="ts">
	import { onMount } from 'svelte';
    export let source = "";

    let audio: HTMLAudioElement;
    
    let currentVolume: number;
    $: currentVolume = 50;

    let muted: boolean;
    $: muted = audio?.muted;

    let playing: boolean;
    $: playing = false;

    let progress: number;
    $: progress = 0;

    let maxProgress: number;
    $: maxProgress = 0;

    const playPause = () => {
        if (audio.paused || audio.ended) {
            audio.play();
        } else {
            audio.pause();
        }

        playing = !(audio.paused || audio.ended);
    }

    const volumeUpdate = () => {
        if (currentVolume === 0) {
            muted = audio.muted = true;
        } else {
            muted = audio.muted = false;
        }

        audio.volume = currentVolume / 100;
    }

    const muteUnmute = () => {
        muted = audio.muted = !muted;
        currentVolume = muted ? 0 : audio.volume * 100;
    }

	onMount(async () => {
        audio.addEventListener('loadedmetadata', () => {
            maxProgress = audio.duration;
            volumeUpdate();
        });
        audio.addEventListener('timeupdate', () => {
            progress = audio.currentTime;
        });
	});
</script>

<div class="player">
    <audio loop preload="auto" src="{source}" bind:this={audio}></audio>
    <progress max="{maxProgress}" value={progress}>
        <div class="progress-bar"></div>
    </progress>
    <div class="controls">
        <button type="button" on:click={playPause}>
            {#if playing}
                <i class="fa-solid fa-pause"></i>
            {:else}
                <i class="fa-solid fa-play"></i>
            {/if}
        </button>
        <button type="button" on:click={muteUnmute}>
            {#if currentVolume === 0 || muted}
                <i class="fa-solid fa-volume-off"></i>
            {:else if currentVolume < 50}
                <i class="fa-solid fa-volume-low"></i>
            {:else}
                <i class="fa-solid fa-volume-high"></i>
            {/if}        
        </button>
        <input type="range" name="volume" min="0" max="100" bind:value={currentVolume} on:change={volumeUpdate} />
    </div>
</div>

<style lang="scss">
    progress {
        width: 100%;
    }
    .controls {
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: center;
        gap: 1rem;
    }
</style>
<script lang="ts">
	import { onMount } from 'svelte';
	import { setAudioProgress } from '$lib/store/audioProgress';
	export let source = '';

	let audio: HTMLAudioElement;

	let currentVolume: number;
	$: currentVolume = 100;

	let muted: boolean;
	$: muted = audio?.muted;

	let isPlaying: boolean;
	$: isPlaying = false;

	let maxProgress: number;
	$: maxProgress = 0;

	const playPause = () => {
		if (!audio) return;

		if (audio.paused || audio.ended) {
			play();
		} else {
			pause();
		}

		isPlaying = !(audio.paused || audio.ended);
	};

	const volumeUpdate = () => {
		if (!audio) return;

		if (currentVolume === 0) {
			muted = audio.muted = true;
		} else {
			muted = audio.muted = false;
		}

		audio.volume = currentVolume / 100;
	};

	const muteUnmute = () => {
		if (!audio) return;

		muted = audio.muted = !muted;
		currentVolume = muted ? 0 : audio.volume * 100;
	};

	const updatePlayState = () => {
		isPlaying = !(audio.paused || audio.ended);
		setAudioProgress({ isPlaying });
	};

	export const play = () => {
		if (!audio) return;

		audio.play();

		updatePlayState();
	};

	export const pause = () => {
		if (!audio) return;

		audio.pause();

		updatePlayState();
	};

	onMount(async () => {
		isPlaying = false;

		audio.addEventListener('loadedmetadata', () => {
			setAudioProgress({ progress: 0, max: audio?.duration });
			volumeUpdate();
		});
		audio.addEventListener('timeupdate', () => {
			setAudioProgress({ progress: audio?.currentTime });
		});

		// Loop in JavaScript rather than HTML so that we still get the
		// completed timeupdate event.
		audio.addEventListener('ended', play);
	});
</script>

<div class="player">
	<audio preload="auto" src={source} bind:this={audio}></audio>
	<div class="controls">
		<button type="button" on:click={playPause}>
			{#if isPlaying}
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
		<input
			type="range"
			name="volume"
			min="0"
			max="100"
			class="volume"
			bind:value={currentVolume}
			on:change={volumeUpdate}
		/>
	</div>
</div>

<style lang="scss">
	.controls {
		display: flex;
		flex-direction: row;
		align-items: center;
		justify-content: center;
		gap: 1rem;

		.volume {
			display: none;

			@media (pointer: fine) {
				display: block;
			}
		}
	}
</style>

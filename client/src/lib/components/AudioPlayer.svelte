<script lang="ts">
	import { onMount } from 'svelte';
	import { setAudioProgress } from '$lib/store/audioProgress';
	export let source = '';
	export let disabled = false;

	let audio: HTMLAudioElement;

	let currentVolume: number;
	$: currentVolume = 100;

	let muted: boolean;
	$: muted = false;

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

	const volumeUpdate = (event: Event) => {
		if (!audio) return;

		const target = event?.currentTarget as HTMLInputElement;
		const volume = parseInt(target?.value ?? '0');

		if (volume === 0) {
			muted = audio.muted = true;
		} else {
			muted = audio.muted = false;
		}

		setVolume(volume);
	};

	const setVolume = (volume: number) => {
		if (!audio) return;

		currentVolume = volume;
		audio.volume = volume / 100;
		localStorage.setItem('volume', currentVolume.toString());
	};

	const toggleMute = () => setMuted(!muted);

	const setMuted = (toMute: boolean) => {
		if (!audio) return;

		muted = audio.muted = toMute;
		currentVolume = muted ? 0 : audio.volume * 100;
		localStorage.setItem('muted', muted.toString());
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
		currentVolume = parseInt(localStorage.getItem('volume') ?? '100');
		muted = localStorage.getItem('muted') === 'true';

		audio.addEventListener('loadedmetadata', () => {
			setAudioProgress({ progress: 0, max: audio?.duration });
			setVolume(currentVolume);
			setMuted(muted);
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
		<button type="button" {disabled} on:click={playPause}>
			{#if isPlaying}
				<i class="fa-solid fa-pause"></i>
			{:else}
				<i class="fa-solid fa-play"></i>
			{/if}
		</button>
		<div class="volume-controls">
			<button type="button" on:click={toggleMute}>
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
</div>

<style lang="scss">
	.controls {
		--gap: 2.5rem;

		display: flex;
		flex-direction: row;
		align-items: center;
		justify-content: center;
		gap: var(--gap);
		font-size: 1.5rem;

		button {
			display: flex;
			align-items: center;

			i {
				width: 1em;
				aspect-ratio: 1;
			}
		}

		.volume-controls {
			display: flex;
			flex-direction: row;
			gap: 1rem;
			align-items: center;

			.volume {
				width: 0px;
			}
		}

		@media (pointer: fine) {
			.volume-controls:hover {
				.volume {
					width: 100%;
					transition-property: width;
					transition-timing-function: var(--animation-timing);
					transition-duration: var(--animation-speed-quick);
					transition-delay: 200ms;
				}
			}
		}
	}
</style>

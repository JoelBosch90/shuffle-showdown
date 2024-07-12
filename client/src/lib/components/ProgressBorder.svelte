<script lang="ts">
	import { progressState } from '$lib/store/progress';
	import type { ProgressState } from '$lib/types/ProgressState';

	let container: HTMLDivElement;
	let currentUpdate = Promise.resolve();
	let topProgressInPixels = 0;
	let rightProgressInPixels = 0;
	let bottomProgressInPixels = 0;
	let leftProgressInPixels = 0;
	let progressSpeed = 0;

	const timeout = async (ms: number) => new Promise<void>((resolve) => setTimeout(resolve, ms));

	const update = async ({ progress, max }: ProgressState) => {
		if (!container) return;

		const { clientHeight, clientWidth } = container;
		const totalDistance = 2 * (clientHeight + clientWidth);
		const newProgressFraction = progress / Math.max(progress, max);
		const newProgressInPixels = totalDistance * newProgressFraction;

		const oldProgressInPixels =
			topProgressInPixels + rightProgressInPixels + bottomProgressInPixels + leftProgressInPixels;
		const pixelsToTravel = newProgressInPixels - oldProgressInPixels;
		const timeToTravel = (pixelsToTravel / totalDistance) * max;
		const isReset = timeToTravel < 0;

		await timeout(timeToTravel);

		progressSpeed = isReset ? 0 : timeToTravel;

		topProgressInPixels = Math.min(clientWidth, newProgressInPixels);
		rightProgressInPixels = Math.min(clientHeight, newProgressInPixels - topProgressInPixels);
		bottomProgressInPixels = Math.min(
			clientWidth,
			newProgressInPixels - topProgressInPixels - rightProgressInPixels
		);
		leftProgressInPixels = Math.min(
			clientHeight,
			newProgressInPixels - topProgressInPixels - rightProgressInPixels - bottomProgressInPixels
		);
	};

	progressState.subscribe((newState) => {
		currentUpdate = currentUpdate.then(() => update(newState));
	});
</script>

<div
	class="progress-background progressing"
	style="--top-progress: {topProgressInPixels}px; --right-progress: {rightProgressInPixels}px; --bottom-progress: {bottomProgressInPixels}px; --left-progress: {leftProgressInPixels}px; --progress-speed: {progressSpeed}s;"
	bind:this={container}
>
	<div class="progress-slider-top"></div>
	<div class="progress-slider-left"></div>
	<div class="progress-slider-bottom"></div>
	<div class="progress-slider-right"></div>
	<div class="content-wrapper"><slot /></div>
</div>

<style lang="scss">
	.progress-background {
		--progress-background-padding: 0.5rem;
		--progress-background-default: var(--white);
		--top-progress: 0px;
		--right-progress: 0px;
		--bottom-progress: 0px;
		--left-progress: 0px;

		$progress-slider: calc(100% - 2 * var(--progress-background-padding));

		display: flex;
		width: 100%;
		padding: var(--progress-background-padding);
		background-color: var(--progress-background-default);

		&.progressing {
			background: var(--rainbow);
			background-size: var(--animation-background-flow-size);
			animation: animation-background-flow 10s linear infinite reverse;
		}

		.progress-slider-top,
		.progress-slider-right,
		.progress-slider-bottom,
		.progress-slider-left {
			content: '';
			position: absolute;
			inset: 0;
			background-color: var(--progress-background-default);
			transition:
				width var(--progress-speed) linear,
				height var(--progress-speed) linear;
		}

		.progress-slider-top {
			bottom: $progress-slider;
			left: auto;
			width: calc(100% - var(--top-progress));
		}

		.progress-slider-right {
			top: auto;
			left: $progress-slider;
			$right-height-available: calc(100% - var(--progress-background-padding));
			height: calc($right-height-available - var(--right-progress));
		}

		.progress-slider-bottom {
			right: var(--progress-background-padding);
			top: $progress-slider;
			$bottom-width-available: calc(100% - var(--progress-background-padding));
			width: calc($bottom-width-available - var(--bottom-progress));
		}

		.progress-slider-left {
			inset: var(--progress-background-padding) $progress-slider var(--progress-background-padding)
				0;
			$left-height-available: calc(100% - 2 * var(--progress-background-padding));
			height: calc($left-height-available - var(--left-progress));
		}

		.content-wrapper {
			position: absolute;
			inset: var(--progress-background-padding);
			border-radius: var(--progress-background-padding);
			display: flex;
			flex-direction: row;
			justify-content: center;
			background-color: var(--progress-background-default);
		}
	}
</style>

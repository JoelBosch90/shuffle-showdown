<script lang="ts">
	import { audioProgressState } from '$lib/store/audioProgress';
	import type { AudioProgressState } from '$lib/types/AudioProgressState';
	import { wait } from '$lib/helpers/wait';

	type ProgressPerSide = {
		top: number;
		right: number;
		bottom: number;
		left: number;
		total: number;
		maxDistance: number;
	};

	const MILLISECONDS_IN_SECOND = 1000;
	const DEFAULT_PROGRESS: ProgressPerSide = {
		top: 0,
		right: 0,
		bottom: 0,
		left: 0,
		total: 0,
		maxDistance: 0
	};
	const BORDER_WIDTH = 10;

	let container: HTMLDivElement;
	let currentUpdate = Promise.resolve();
	let progressSpeed = 0;
	let isProgressing: boolean = false;
	const currentProgress: ProgressPerSide = DEFAULT_PROGRESS;

	const calculateTotalDistance = (clientHeight: number, clientWidth: number) =>
		2 * (clientHeight + clientWidth);

	const calculateSideProgress = (progress: number, sideMax: number, previousProgress: number) =>
		Math.max(0, Math.min(sideMax, progress - previousProgress));

	const calculateProgress = (progress: number, max: number): ProgressPerSide => {
		if (!container) return DEFAULT_PROGRESS;

		const { clientHeight, clientWidth } = container;
		const maxDistance = calculateTotalDistance(clientHeight, clientWidth);
		const progressFraction = progress / Math.max(progress, max);
		const total = maxDistance * progressFraction;

		let accumulated = 0;
		const top = calculateSideProgress(total, clientWidth, accumulated);
		accumulated += top;

		const right = calculateSideProgress(total, clientHeight - BORDER_WIDTH, accumulated);
		accumulated += right;

		const bottom = calculateSideProgress(total, clientWidth - BORDER_WIDTH, accumulated);
		accumulated += bottom;

		const left = calculateSideProgress(total, clientHeight - 2 * BORDER_WIDTH, accumulated);

		return { top, right, bottom, left, total, maxDistance };
	};

	const calculateAnimationTime = ({ total, maxDistance }: ProgressPerSide, max: number) => {
		const oldTotal = currentProgress.total;
		const pixelsToTravel = total - oldTotal;
		const animationTime = (pixelsToTravel / maxDistance) * max;

		return total === 0 ? 0 : animationTime;
	};

	const animateProgress = async (progressPerSide: ProgressPerSide, max: number) => {
		const animationTime = calculateAnimationTime(progressPerSide, max);
		const { top, right, bottom, left, total, maxDistance } = progressPerSide;

		progressSpeed = animationTime;
		currentProgress.top = top;
		currentProgress.right = right;
		currentProgress.bottom = bottom;
		currentProgress.left = left;
		currentProgress.total = total;
		currentProgress.maxDistance = maxDistance;

		await wait(animationTime);
	};

	const update = async ({ progress, max, isPlaying }: AudioProgressState) => {
		isProgressing = isPlaying;
		await animateProgress(calculateProgress(progress, max), max);
	};

	audioProgressState.subscribe((newState) => {
		currentUpdate = currentUpdate.then(() => update(newState));
	});
</script>

<div
	class="progress-background"
	style="
    --border-width: {BORDER_WIDTH}px;
    --top-progress: {currentProgress.top}px; 
    --right-progress: {currentProgress.right}px; 
    --bottom-progress: {currentProgress.bottom}px; 
    --left-progress: {currentProgress.left}px; 
    --progress-speed: {progressSpeed * MILLISECONDS_IN_SECOND}ms;
    --play-state: {isProgressing ? 'running' : 'paused'}
  "
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
		--progress-background-default: var(--white);
		--border-width: 0px;
		--top-progress: 0px;
		--right-progress: 0px;
		--bottom-progress: 0px;
		--left-progress: 0px;
		--progress-speed: 0s;

		$progress-slider: calc(100% - 2 * var(--border-width));

		display: flex;
		width: 100%;
		padding: var(--border-width);
		background-color: var(--progress-background-default);

		background: var(--rainbow);
		background-size: var(--animation-background-flow-size);
		animation: animation-background-flow 10s linear infinite reverse;
		animation-play-state: var(--play-state);

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
			$right-height-available: calc(100% - var(--border-width));
			height: calc($right-height-available - var(--right-progress));
		}

		.progress-slider-bottom {
			right: var(--border-width);
			top: $progress-slider;
			$bottom-width-available: calc(100% - var(--border-width));
			width: calc($bottom-width-available - var(--bottom-progress));
		}

		.progress-slider-left {
			inset: var(--border-width) $progress-slider var(--border-width) 0;
			$left-height-available: calc(100% - 2 * var(--border-width));
			height: calc($left-height-available - var(--left-progress));
		}

		.content-wrapper {
			position: absolute;
			inset: var(--border-width);
			display: flex;
			flex-direction: row;
			justify-content: center;
			border-radius: var(--border-width);
			padding: 0.25rem;
			box-sizing: border-box;
			background-color: var(--progress-background-default);
		}
	}
</style>

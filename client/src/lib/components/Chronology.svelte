<script lang="ts">
	import type { WonTrack } from '$lib/types/WonTrack';
	import type { Artist } from '$lib/types/Artist';
	import type { Answer } from '$lib/types/Answer';
	import { onMount } from 'svelte';
	import AudioPlayer from './AudioPlayer.svelte';

	export let disabled: boolean = false;
	export let wonTracks: WonTrack[] = [];
	export let onSelect: (answer: Answer) => void;

	interface Card {
		releaseYear: string;
		name?: string;
		artists?: string;
		isGuess?: boolean;
	}

	let container: HTMLOListElement;

	const RATIO_OF_SCREEN_FOR_FULL_MOVE = 0.8;
	const INDEX_CHANGE_FOR_FULL_MOVE = 4;
	const UNKNOWN_RELEASE_YEAR = '???';
	const guessCard: Card = {
		releaseYear: UNKNOWN_RELEASE_YEAR,
		isGuess: true
	};

	const joinArtists = (artists: Artist[] = []) => {
		return artists.map((artist) => artist.name).join(', ');
	};

	const selectAnswer = () => {
		if (disabled) return;

		const cardBefore = guessIndex > 0 ? cards[guessIndex - 1] : undefined;
		const cardAfter = guessIndex < cards.length - 1 ? cards[guessIndex + 1] : undefined;

		onSelect({
			afterReleaseYear: cardBefore ? parseInt(cardBefore?.releaseYear ?? '') : undefined,
			beforeReleaseYear: cardAfter ? parseInt(cardAfter?.releaseYear ?? '') : undefined,
			guessIndex: guessIndex
		});
	};

	const trackToCard = ({ track }: WonTrack): Card => ({
		releaseYear: track.releaseYear?.toString() ?? UNKNOWN_RELEASE_YEAR,
		name: track.name,
		artists: joinArtists(track.artists)
	});

	const sortWonTracks = (a: WonTrack, b: WonTrack) => {
		const aReleaseYear = a.track.releaseYear ?? 0;
		const bReleaseYear = b.track.releaseYear ?? 0;
		return aReleaseYear - bReleaseYear;
	};

	const areNewTracks = (tracks: WonTrack[]) => {
		if (tracks.length !== wonTracks.length) return true;

		const cards = [...trackCards];
		const hasNewTrack = tracks.some((track) =>
			cards.every((card) => card.name !== track.track.name)
		);
		const hasLostTrack = cards.some((card) =>
			tracks.every((track) => track.track.name !== card.name)
		);

		return hasNewTrack || hasLostTrack;
	};

	const onWheelEvent = (event: WheelEvent) => {
		if (disabled) return;

		const direction = event.deltaY + event.deltaX > 0 ? 1 : -1;

		guessIndex = clampGuessIndex(guessIndex + direction);
	};

	const clampGuessIndex = (newGuessIndex: number) => {
		return Math.max(0, Math.min(newGuessIndex, trackCards.length));
	};

	const getClientLocation = (event: MouseEvent | TouchEvent) => {
		if (event instanceof MouseEvent) return { x: event.clientX, y: event.clientY };
		return { x: event.touches[0].clientX, y: event.touches[0].clientY };
	};

	const onMoveStart = (event: MouseEvent | TouchEvent) => {
		if (disabled) return;

		const { x: startX, y: startY } = getClientLocation(event);
		const startGuessIndex = guessIndex;

		const onMove = (event: MouseEvent | TouchEvent) => {
			const { x: currentX, y: currentY } = getClientLocation(event);
			const heightPercentageTraveled =
				(currentX - startX) / (container.clientWidth * RATIO_OF_SCREEN_FOR_FULL_MOVE);
			const widthPercentageTraveled =
				(currentY - startY) / (container.clientHeight * RATIO_OF_SCREEN_FOR_FULL_MOVE);

			const distance = Math.hypot(widthPercentageTraveled, heightPercentageTraveled);
			const indexChange = Math.floor(distance * INDEX_CHANGE_FOR_FULL_MOVE);
			const direction = widthPercentageTraveled - heightPercentageTraveled > 0 ? 1 : -1;

			guessIndex = clampGuessIndex(startGuessIndex + direction * indexChange);
		};

		const onMoveEnd = () => {
			window.removeEventListener('touchmove', onMove);
			window.removeEventListener('touchend', onMoveEnd);
			window.removeEventListener('touchcancel', onMoveEnd);

			window.removeEventListener('mousemove', onMove);
			window.removeEventListener('mouseup', onMoveEnd);
		};

		window.addEventListener('touchmove', onMove);
		window.addEventListener('touchend', onMoveEnd);
		window.addEventListener('touchcancel', onMoveEnd);

		window.addEventListener('mousemove', onMove);
		window.addEventListener('mouseup', onMoveEnd);
	};

	export const selectIndex = (newIndex: number) => {
		if (!disabled) return;
		if (typeof newIndex !== 'number') return;

		guessIndex = newIndex;
	};

	$: wonTracks,
		(() => {
			if (!areNewTracks(wonTracks)) return;
			trackCards = wonTracks.toSorted(sortWonTracks).map(trackToCard);
			guessIndex = Math.ceil(trackCards.length / 2);
		})();

	let trackCards: Card[] = [];
	let guessIndex: number = 0;

	let cards: Card[];
	$: cards = [...trackCards?.slice(0, guessIndex), guessCard, ...trackCards?.slice(guessIndex)];
	$: cards, selectAnswer();

	onMount(() => {
		container.addEventListener('wheel', onWheelEvent);
		container.addEventListener('mousedown', onMoveStart);
		container.addEventListener('touchstart', onMoveStart);
	});
</script>

<ol class="chronology" class:disabled bind:this={container}>
	{#each cards as card, cardIndex}
		<li
			class="card {card.isGuess ? 'guess' : ''}"
			style="--normalized-index: {cardIndex - guessIndex}"
		>
			<h2>{card.releaseYear}</h2>

			{#if card.name}
				<p class="track">{card.name}</p>
			{/if}

			{#if card.artists}
				<p class="artist">{card.artists}</p>
			{/if}
		</li>
	{/each}
</ol>

<style lang="scss">
	.chronology {
		container-type: size;
		container-name: chronology;
		display: flex;
		flex-grow: 1;
		width: 100%;
		margin: 0;
		padding: 0;
		overflow: hidden;
		position: relative;
		box-sizing: border-box;

		cursor: grab;
		user-select: none;
		-webkit-user-select: none;

		&.disabled {
			cursor: not-allowed;
			.card {
				color: var(--gray-dark);
			}
		}

		.card {
			--normalized-index: 0;
			--card-padding: 1rem;
			--card-font-size: 12cqw;
			--card-min-height: 6rem;
			--card-scale-height: 50cqmin;
			--card-max-height: 16rem;
			--card-aspect-ratio: 1 / 1.25;
			--card-box-shadow-space: 0.75rem;
			--center-card-margin-percentage: 0.33;

			$normalized-index: var(--normalized-index);
			$distance-from-guess-card: max($normalized-index, -1 * $normalized-index);
			$card-height: clamp(var(--card-min-height), var(--card-scale-height), var(--card-max-height));
			$card-width: calc($card-height * var(--card-aspect-ratio));

			display: flex;
			box-sizing: border-box;
			container-type: size;
			container-name: card;

			position: absolute;
			top: 0;
			left: 0;
			aspect-ratio: var(--card-aspect-ratio);
			height: $card-height;
			z-index: calc(var(--card-level) - $distance-from-guess-card);

			flex-direction: column;
			justify-content: center;
			align-items: center;
			gap: 0.2rem;

			list-style: none;
			padding: var(--card-padding);
			border-radius: var(--card-padding);
			box-shadow: 0 0 var(--card-box-shadow-space) rgba(0, 0, 0, 0.25);
			background-color: var(--white);
			overflow: hidden;

			transition:
				transform var(--animation-speed-quick),
				z-index var(--animation-speed-quick),
				opacity var(--animation-speed-quick),
				border var(--animation-speed-quick),
				color var(--animation-speed-quick);

			$direction: clamp(-1, $normalized-index, 1);
			$non-zero-index: calc($distance-from-guess-card + var(--center-card-margin-percentage));
			$fraction: calc(1 / $non-zero-index);
			$clamped-fraction: clamp(0, $fraction, 1);
			$scale: clamp(0, $clamped-fraction, 1);
			$index-offset: calc($scale - 1);
			$card-offset: calc($index-offset * 50 * $direction);
			$size-scale: calc($scale * 0.5 + 0.5);

			@function calculated-offset($side-unit, $card-side-length) {
				@return calc(
					$card-offset * $side-unit - $direction * $card-side-length *
						var(--center-card-margin-percentage)
				);
			}
			$card-vertical-offset: calculated-offset(1cqh, $card-height);
			$card-horizontal-offset: calculated-offset(1cqw, $card-width);

			@function scaled-value($value) {
				@return calc($value * $size-scale);
			}
			$card-scaled-height: scaled-value($card-height);
			$card-scaled-width: scaled-value($card-width);
			$card-scaled-box-shadow-space: scaled-value(var(--card-box-shadow-space));

			@function reverse-scaled-value($value) {
				@return calc($value * (1 - $size-scale));
			}
			$card-scaling-height-loss: reverse-scaled-value($card-height);
			$card-scaling-width-loss: reverse-scaled-value($card-width);

			@function start-value($card-side-length) {
				@return calc(
					0% - reverse-scaled-value($card-side-length) * 0.5 + $card-scaled-box-shadow-space
				);
			}
			$card-at-top: start-value($card-height);
			$card-at-left: start-value($card-width);

			@function middle-value($side-unit, $card-side-length) {
				@return calc(
					50 * $side-unit -
						(scaled-value($card-side-length) + reverse-scaled-value($card-side-length)) * 0.5
				);
			}
			$card-at-vertical-center: middle-value(1cqh, $card-height);
			$card-at-horizontal-center: middle-value(1cqw, $card-width);

			@function end-value($side-unit, $card-side-length) {
				@return calc(
					100 * $side-unit -
						(
							$card-side-length - reverse-scaled-value($card-side-length) * 0.5 +
								$card-scaled-box-shadow-space
						)
				);
			}
			$card-at-bottom: end-value(1cqh, $card-height);
			$card-at-right: end-value(1cqw, $card-width);

			$vertical-transform: translateY(
				clamp($card-at-top, calc($card-at-vertical-center + $card-vertical-offset), $card-at-bottom)
			);
			$horizontal-transform: translateX(
				clamp(
					$card-at-left,
					calc($card-at-horizontal-center - $card-horizontal-offset),
					$card-at-right
				)
			);
			transform: $vertical-transform $horizontal-transform scale($size-scale);

			p,
			h2 {
				max-width: 100%;
				margin: 0;
				overflow: hidden;
				text-overflow: ellipsis;
				text-align: center;
			}

			p {
				display: -webkit-box;
				-webkit-line-clamp: 2;
				-webkit-box-orient: vertical;
				font-size: var(--card-font-size);

				&.track {
					font-size: calc(var(--card-font-size) * 1.5);
				}
			}

			h2 {
				font-size: calc(var(--card-font-size) * 3.75);
			}

			&.guess {
				opacity: 50%;
				border: 2px dashed var(--gray-dark);
			}
		}
	}
</style>

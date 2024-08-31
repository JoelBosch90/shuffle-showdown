<script lang="ts">
	import SongCard from '$lib/components/SongCard.svelte';
	import type { WonTrack } from '$lib/types/WonTrack';
	import type { Artist } from '$lib/types/Artist';
	import type { Answer } from '$lib/types/Answer';
	import { onMount } from 'svelte';

	export let cardColor: string = '';
	export let cardIcon: string = '';
	export let disabled: boolean = false;
	export let wonTracks: WonTrack[] = [];
	export let onSelect: (answer: Answer) => void;

	interface Card {
		releaseYear: string;
		name?: string;
		artists?: Artist[];
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
		artists: track.artists
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
	let beforeTextIndex: number;
	let afterTextIndex: number;
	let isBeforeTextShown: boolean;
	let isAfterTextShown: boolean;

	let cards: Card[];
	$: cards = [...trackCards?.slice(0, guessIndex), guessCard, ...trackCards?.slice(guessIndex)];
	$: cards, selectAnswer();
	$: cards, (beforeTextIndex = -1 * guessIndex - 1), (isBeforeTextShown = beforeTextIndex === -1);
	$: cards,
		(afterTextIndex = cards?.length - guessIndex),
		(isAfterTextShown = afterTextIndex === 1);

	onMount(() => {
		container.addEventListener('wheel', onWheelEvent);
		container.addEventListener('mousedown', onMoveStart);
		container.addEventListener('touchstart', onMoveStart);
	});
</script>

<ol class="chronology" class:disabled bind:this={container}>
	{#each cards as card, cardIndex}
		<SongCard
			{cardColor}
			{cardIcon}
			{disabled}
			normalizedIndex={cardIndex - guessIndex}
			isGuess={card.isGuess}
			releaseYear={card.releaseYear}
			trackName={card.name}
			artists={card.artists}
		/>
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
		}
	}
</style>

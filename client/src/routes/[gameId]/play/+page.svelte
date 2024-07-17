<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import AudioPlayer from '$lib/components/AudioPlayer.svelte';
	import Chronology from '$lib/components/Chronology.svelte';
	import Celebration from '$lib/components/Celebration.svelte';
	import type { Player } from '$lib/types/Player';
	import type { Round } from '$lib/types/Round';
	import type { GameSessionUpdate } from '$lib/types/GameSessionUpdate';
	import type { Answer } from '$lib/types/Answer';
	import { GameSession } from '$lib/services/GameSession';
	import { findPlayerInGameSessionUpdate } from '$lib/helpers/findPlayerInGameSessionUpdate';
	import LoadingButton from '$lib/components/LoadingButton.svelte';
	import { debounce } from '$lib/helpers/debounce';

	const DEBOUNCE_WAIT_MILLISECONDS = 150;
	const gameId = $page.params.gameId;
	let session: GameSession | void | null = null;

	let gameUpdate: GameSessionUpdate | null;
	$: gameUpdate = null;

	let currentRound: Round | null;
	$: currentRound = null;

	let currentPlayer: Player | null;
	$: currentPlayer = null;

	let me: Player | null;
	$: me = null;

	let isPlaying: boolean;
	$: isPlaying = false;

	let isLoading = false;

	let selectedAnswer: Answer | null = null;

	let audioPlayer: AudioPlayer | null = null;

	let chronology: Chronology | null = null;

	let celebration: Celebration | null = null;

	const getCurrentRound = (update: GameSessionUpdate | null): Round | null => {
		if (!update) return null;

		const maxRoundNumber = Math.max(...update.rounds.map((round) => round.number));
		const currentRound = update.rounds.find((round) => round.number === maxRoundNumber);

		return currentRound ?? null;
	};

	const debouncedAnswerSelectionUpdate = debounce(
		(answer) => session?.updateAnswerSelection(answer as Answer),
		DEBOUNCE_WAIT_MILLISECONDS
	);

	const onAnswerSelect = (answer: Answer) => {
		selectedAnswer = answer;
		debouncedAnswerSelectionUpdate(answer);
	};

	const onAnswerSubmit = () => {
		if (!selectedAnswer) return;
		isLoading = true;
		session?.submitAnswer(selectedAnswer);
		audioPlayer?.pause();
	};

	const celebrate = ({
		game: update,
		me: newMe
	}: {
		game: GameSessionUpdate | null;
		me: Player | null;
	}) => {
		celebration?.update({
			oldUpdate: gameUpdate,
			newUpdate: update,
			oldMe: me,
			newMe
		});
	};

	const updatePage = async ({
		game: update,
		me: newMe
	}: {
		game: GameSessionUpdate | null;
		me: Player | null;
	}) => {
		gameUpdate = update;
		me = newMe;
		currentRound = getCurrentRound(update);
		currentPlayer = findPlayerInGameSessionUpdate(update, currentRound?.playerId);
		isPlaying = !!currentPlayer && currentPlayer.id === me?.id;
		isLoading = false;

		if (!update?.hasStarted) return goto(`/${gameId}/lobby`);
	};

	onMount(async () => {
		if (!session) session = new GameSession(gameId);
		session.onUpdate((gameUpdate) => {
			celebrate(gameUpdate);
			updatePage(gameUpdate);
		});
		session.onAnswerSelectionUpdate(({ answerSelectionUpdate }) => {
			if (
				answerSelectionUpdate?.answer?.guessIndex === undefined ||
				answerSelectionUpdate.playerId === me?.id ||
				isPlaying
			) {
				return;
			}
			chronology?.selectIndex(answerSelectionUpdate.answer.guessIndex);
		});

		const latestUpdate = session.getCachedUpdate();
		if (latestUpdate) updatePage(latestUpdate);

		await session.initialize();
	});
</script>

<svelte:head>
	<title>Game</title>
	<meta
		name="description"
		content="Shuffle Showdown game. First to collect {gameUpdate?.songsToWin} songs by placing them in order of release wins!"
	/>
</svelte:head>

<div class="container">
	<div class="card-field">
		<svelte:component
			this={Chronology}
			bind:this={chronology}
			wonTracks={currentPlayer?.wonTracks}
			onSelect={onAnswerSelect}
			disabled={!isPlaying}
		/>

		<h2 class="round-info">Round {currentRound?.number}</h2>

		{#if currentPlayer}
			<p class="player-info">
				Now playing: {isPlaying ? 'you' : currentPlayer.name} ({currentPlayer.wonTracks
					?.length}/{gameUpdate?.songsToWin})
			</p>
		{/if}
	</div>

	<div class="controls">
		<svelte:component
			this={AudioPlayer}
			disabled={!isPlaying}
			bind:this={audioPlayer}
			source={currentRound?.track.previewUrl}
		/>

		<LoadingButton {isLoading} onClick={onAnswerSubmit} isDisabled={!isPlaying}>
			Select answer
		</LoadingButton>
	</div>
	<svelte:component this={Celebration} bind:this={celebration} />
</div>

<style lang="scss">
	.container {
		--margin-small: 0.25rem;
		--margin-large: 1rem;

		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 100%;
		width: 100%;
		position: relative;

		.card-field {
			position: relative;
			flex: 1;
			overflow: hidden;
			display: flex;
			width: 100%;

			.round-info,
			.player-info {
				position: absolute;
				pointer-events: none;
			}

			.round-info {
				top: 0;
				left: 0;
				margin: var(--margin-small) 0 0 var(--margin-small);
			}

			.player-info {
				right: 0;
				bottom: 0;
				margin: 0 var(--margin-small) var(--margin-small) 0;
			}
		}

		.controls {
			width: 100%;
			display: flex;
			flex-direction: row;
			gap: 1rem;
			align-items: center;
			justify-content: space-between;
			margin-top: var(--margin-small);
		}
	}
</style>

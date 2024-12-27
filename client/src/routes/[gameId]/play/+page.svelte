<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import AudioPlayer from '$lib/components/AudioPlayer.svelte';
	import Chronology from '$lib/components/Chronology.svelte';
	import type { Player } from '$lib/types/Player';
	import type { PlayerWithIconAndColor } from '$lib/types/PlayerWithIconAndColor';
	import type { Round } from '$lib/types/Round';
	import type { GameSessionUpdate } from '$lib/types/GameSessionUpdate';
	import type { Answer } from '$lib/types/Answer';
	import { GameSession } from '$lib/services/GameSession';
	import { findPlayerInGameSessionUpdate } from '$lib/helpers/findPlayerInGameSessionUpdate';
	import LoadingButton from '$lib/components/LoadingButton.svelte';
	import { debounce } from '$lib/helpers/debounce';
	import { showToast } from '$lib/store/toasts';
	import { ToastType } from '$lib/enums/ToastType';
	import { getPlayerColor } from '$lib/helpers/getPlayerColor';
	import { getPlayerIcon } from '$lib/helpers/getPlayerIcon';
	import { wait } from '$lib/helpers/wait';

	const DEBOUNCE_WAIT_MILLISECONDS = 150;
	const gameId = $page.params.gameId;
	const updatePromise = Promise.resolve();

	let session: GameSession | void | null = null;

	let gameUpdate: GameSessionUpdate | null;
	$: gameUpdate = null;

	let currentRound: Round | null;
	$: currentRound = null;

	let currentPlayer: PlayerWithIconAndColor | null;
	$: currentPlayer = null;

	let me: Player | null;
	$: me = null;

	let isPlaying: boolean;
	$: isPlaying = false;

	let isLoading = false;

	let selectedAnswer: Answer | null = null;

	let audioPlayer: AudioPlayer | null = null;

	let chronology: Chronology | null = null;

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

	const revealAnswer = (update: GameSessionUpdate | null) => {
		if (gameUpdate?.rounds.length === update?.rounds.length) {
			return;
		}

		const trackToReveal = update?.rounds.find(
			(round) => round.number === currentRound?.number
		)?.track;
		const trackWon = update?.players
			?.find((player) => player.id === currentRound?.playerId)
			?.wonTracks?.find((wonTrack) => wonTrack.track.name === trackToReveal?.name);

		if (!trackToReveal) return;

		chronology?.reveal({
			releaseYear: trackToReveal.releaseYear?.toString() ?? '???',
			name: trackToReveal.name ?? '',
			artists: trackToReveal.artists ?? [],
			isWon: !!trackWon
		});
	};

	const updatePage = async ({
		game: update,
		me: newMe
	}: {
		game: GameSessionUpdate | null;
		me: Player | null;
	}) => {
		const player = findPlayerInGameSessionUpdate(update, currentRound?.playerId);
		gameUpdate = update;
		me = newMe;
		currentRound = getCurrentRound(update);
		currentPlayer = player
			? {
					...player,
					color: getPlayerColor(player.id),
					icon: getPlayerIcon(player.id)
				}
			: null;
		isPlaying = !!currentPlayer && currentPlayer.id === me?.id;
		isLoading = false;

		if (!update?.hasStarted) return goto(`/${gameId}/lobby`);
	};

	onMount(async () => {
		if (!session) session = new GameSession(gameId);
		session.onUpdate(async (gameUpdate) => {
			updatePromise.then(async () => {
				revealAnswer(gameUpdate.game);
				await wait(5000);
				updatePage(gameUpdate);
			});
		});
		session.onPlaylistLoadingUpdate(({ playlistLoadingUpdate }) => {
			if (!playlistLoadingUpdate?.liveUpdate) return;

			if (playlistLoadingUpdate?.finishedLoading && playlistLoadingUpdate?.tracksVerified === 0) {
				showToast({
					message: `Loaded ${playlistLoadingUpdate?.tracksLoaded} out of ${playlistLoadingUpdate?.tracksTotal} total tracks`,
					type: ToastType.Success
				});
			}

			if (playlistLoadingUpdate?.finishedChecking) {
				showToast({
					message: `Verified ${playlistLoadingUpdate?.tracksVerified} out of ${playlistLoadingUpdate?.tracksLoaded} loaded tracks`,
					type: ToastType.Success
				});
			}
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
			currentRoundNumber={currentRound?.number}
			cardColor={`--player-${currentPlayer?.color}`}
			cardIcon={currentPlayer?.icon}
			onSelect={onAnswerSelect}
			disabled={!isPlaying}
		/>

		<div class="round-info">
			<h2>
				Round {currentRound?.number}
			</h2>

			{#if currentPlayer}
				<span class="player-info" style="--player-color:var(--player-{currentPlayer?.color})">
					<i class={`fa-solid fa-${currentPlayer?.icon} player-icon`} />
					{currentPlayer?.name}
					({currentPlayer.wonTracks?.length}/{gameUpdate?.songsToWin})
				</span>
			{/if}
		</div>
	</div>

	<div class="controls">
		<svelte:component
			this={AudioPlayer}
			disabled={!isPlaying}
			bind:this={audioPlayer}
			source={currentRound?.track.previewUrl}
		/>

		<LoadingButton
			{isLoading}
			onClick={onAnswerSubmit}
			isDisabled={!isPlaying}
			title="Click to select this answer"
		>
			Select answer
		</LoadingButton>
	</div>
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

			.round-info {
				position: absolute;
				pointer-events: none;
				left: 0;
				top: 0;
				margin: var(--margin-small) 0 0 var(--margin-small);

				h2 {
					margin: 0;
				}

				.player-info {
					i {
						color: var(--player-color, inherit);
					}
				}
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

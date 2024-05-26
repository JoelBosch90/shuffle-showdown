<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import AudioPlayer from '$lib/components/AudioPlayer.svelte';
	import Chronology from '$lib/components/Chronology.svelte';
	import Celebration from '$lib/components/Celebration.svelte';
	import { GameSession } from '$lib/services/GameSession';
	import type { Player } from '$lib/types/Player';
	import type { Round } from '$lib/types/Round';
	import type { GameSessionUpdate } from '$lib/types/GameSessionUpdate';
	import type { Answer } from '$lib/types/Answer';
	import type { Track } from '$lib/types/Track';

	const CELEBRATION_DURATION = 3500;
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

	let selectedAnswer: Answer | null = null;

	let audioPlayer: AudioPlayer | null = null;

	let celebration: Celebration | null = null;
	let isCelebrating: boolean = false;
	let celebrationPromise: Promise<void> = Promise.resolve();

	const getCurrentRound = (update: GameSessionUpdate | null) : Round | null => {
		if (!update) return null;

		const maxRoundNumber = Math.max(...update.rounds.map((round) => round.number));
		const currentRound = update.rounds.find((round) => round.number === maxRoundNumber);

		return currentRound ?? null;
	}

	const findPlayer = (update: GameSessionUpdate | null, playerId: string | undefined) => update?.players.find((player) => player.id === playerId) ?? null;

	const onAnswerSelect = (answer: Answer) => {
		selectedAnswer = answer;
	}

	const onAnswerSubmit = () => {
		if (!selectedAnswer) return;
		session?.submitAnswer(selectedAnswer);
		audioPlayer?.pause();
	}

	const findWinner = (update: GameSessionUpdate, trackPreviewUrl?: string) => {
		if (!trackPreviewUrl) return undefined;

		for (let player of update?.players ?? []) {
			for (let wonTrack of player?.wonTracks ?? []) {
				if (wonTrack.track.previewUrl === trackPreviewUrl) return player;
			}
		}

		return undefined;
	}

	const sleep = async (ms: number) : Promise<void> => {
		return new Promise((resolve) => setTimeout(resolve, ms));
	}

	const celebrateTrack = async (track: Track, player: Player, isOtherPlayer: boolean, hasWon: boolean, isFinalWin: boolean, millisecondsToSleep?: number) => {
		isCelebrating = true;
		celebration?.celebrate(track, player, isOtherPlayer, hasWon, isFinalWin);

		if (millisecondsToSleep) {
			celebrationPromise = celebrationPromise?.then(() => sleep(millisecondsToSleep));
			await celebrationPromise;
			isCelebrating = false;
		}
	}

	const celebrateRound = async (update: GameSessionUpdate, roundIndex: number) => {
		const track = update.rounds.find((round) => round.number === roundIndex)?.track;
		if (!track?.name) return;

		const winner = findWinner(update, track.previewUrl);
		const player = winner ?? findPlayer(update, update.rounds[roundIndex].playerId);

		const isOtherPlayer = player?.id !== me?.id;
		const hasWon = !!winner;

		const wonTracksCount = player?.wonTracks?.length ?? 0;
		const isFinalWin = hasWon && wonTracksCount >= update.songsToWin;

		if (!track || !player) return;
		await celebrateTrack(track, player, isOtherPlayer, hasWon, isFinalWin, CELEBRATION_DURATION);
	}

	const celebrateStart = async (newMe: Player) => {
		if (!newMe) return;

		for (const track of newMe?.wonTracks ?? []) {
			await celebrateTrack(track.track, newMe, false, true, false, CELEBRATION_DURATION);
		}
	}

	const celebrateEnd = async (update: GameSessionUpdate, newMe: Player | null) => {
		if (!update || !update.hasFinished) return;

		const lastRound = getCurrentRound(update);
		if (!lastRound) return;

		const winner = findWinner(update, lastRound.track.previewUrl);
		if (!winner) return;

		await celebrateTrack(lastRound.track, winner, winner.id !== newMe?.id, true, true);
	}

	const processCelebrations = async (currentGame: GameSessionUpdate | null, update: GameSessionUpdate | null, newMe: Player | null) => {
		if (!update) return;

		if (update.hasFinished) await celebrateEnd(update, newMe);

		const lastSeenRound = getCurrentRound(currentGame)?.number ?? 0;
		const currentRound = getCurrentRound(update)?.number ?? 0;
		if (lastSeenRound >= currentRound) return;

		for (let celebratedRound = lastSeenRound; celebratedRound < currentRound; celebratedRound++) {
			celebrateRound(update, celebratedRound);
		}
	}

	const updateState = ({ game: update, me: newMe }: { game: GameSessionUpdate | null, me: Player | null }) => {
		gameUpdate = update;
		me = newMe;
		currentRound = getCurrentRound(update);
		currentPlayer = findPlayer(update, currentRound?.playerId);
		isPlaying = !!currentPlayer && currentPlayer.id === me?.id;
	}

	const onUpdate = async ({ game: update, me: newMe }: { game: GameSessionUpdate | null, me: Player | null }) => {
		processCelebrations(gameUpdate, update, newMe);

		updateState({ game: update, me: newMe });

		if (!update?.hasStarted) return goto(`/game/${gameId}/lobby`);
	}

	onMount(async () => {
		session = new GameSession(gameId);
		session.onUpdate(onUpdate);

		const latestUpdate = session.getCachedUpdate();
		if (latestUpdate) updateState(latestUpdate);
		if (latestUpdate && latestUpdate.game?.rounds.length <= 1) await celebrateStart(latestUpdate.me);
		if (latestUpdate && latestUpdate.game?.hasFinished) await celebrateEnd(latestUpdate.game, latestUpdate.me);

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
	<div class="overlay" class:hidden={!isCelebrating}>
		<svelte:component this={Celebration} bind:this={celebration} />
	</div>
	<div class="game-interface" class:hidden={isCelebrating}>
		<h1>Round {currentRound?.number}</h1>
		{#if currentPlayer}
			<span>Currently playing: {currentPlayer.id === me?.id ? "you" : currentPlayer.name}</span>
			<span>{currentPlayer.id === me?.id ? "You have" : currentPlayer.name + " has"} won {currentPlayer.wonTracks?.length} out of {gameUpdate?.songsToWin} tracks.</span>
		{/if}
		<Chronology wonTracks={currentPlayer?.wonTracks} onSelect={onAnswerSelect} disabled={!isPlaying}/>
		<svelte:component this={AudioPlayer} bind:this={audioPlayer} source="{currentRound?.track.previewUrl}" />
		<button class="filled" on:click={onAnswerSubmit}>
			Select answer
		</button>
	</div>
</div>

<style lang="scss">
	.container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 100%;
		width: 100%;
	}

	.overlay {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		align-items: center;
	}

	.game-interface {
		flex: 1;
		overflow: hidden;
		display: flex;
		flex-direction: column;
		gap: 1rem;
		align-items: center;

		h1 {
			text-align: center;
		}
	}

	.hidden {
		display: none;
	}
</style>

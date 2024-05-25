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
	import type { GameUpdate } from '$lib/types/GameUpdate';
	import type { Answer } from '$lib/types/Answer';
	import type { Track } from '$lib/types/Track';

	const CELEBRATION_DURATION = 3500;
	const gameId = $page.params.gameId;
	let session: GameSession | void | null = null;

	let game: GameUpdate | null;
	$: game = null;

	let currentRound: Round | null;
	$: currentRound = null;
	
	let currentPlayer: Player | null;
	$: currentPlayer = null;

	let me: Player | null;
	$: me = null;

	let isPlaying: boolean;
	$: isPlaying = false;

	let selectedAnswer: Answer | null = null;

	let isCelebrating: boolean = false;
	let celebration: Celebration | null = null;

	const getCurrentRound = (update: GameUpdate | null) : Round | null => {
		if (!update) return null;

		const maxRoundNumber = Math.max(...update.rounds.map((round) => round.number));
		const currentRound = update.rounds.find((round) => round.number === maxRoundNumber);

		return currentRound ?? null;
	}

	const findPlayer = (update: GameUpdate | null, playerId: string | undefined) => update?.players.find((player) => player.id === playerId) ?? null;

	const onAnswerSelect = (answer: Answer) => {
		selectedAnswer = answer;
	}

	const onAnswerSubmit = () => {
		if (!selectedAnswer) return;
		session?.submitAnswer(selectedAnswer);
	}

	const findWinner = (update: GameUpdate, trackPreviewUrl?: string) => {
		if (!trackPreviewUrl) return undefined;

		for (let player of update?.players ?? []) {
			for (let wonTrack of player?.wonTracks ?? []) {
				if (wonTrack.track.previewUrl === trackPreviewUrl) return player;
			}
		}

		return undefined;
	}

	const sleep = async (ms: number) => {
		return new Promise((resolve) => setTimeout(resolve, ms));
	}

	const celebrateTrack = async (track: Track, player: Player, isOtherPlayer: boolean, hasWon: boolean, isFinalWin: boolean, millisecondsToSleep?: number) => {
		isCelebrating = true;
		celebration?.celebrate(track, player, isOtherPlayer, hasWon, isFinalWin);

		if (millisecondsToSleep) {
			await sleep(millisecondsToSleep);
			isCelebrating = false;
		}
	}

	const celebrateRound = async (update: GameUpdate, roundIndex: number) => {
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

	const celebrateEnd = async (update: GameUpdate, newMe: Player | null) => {
		if (!update || !update.hasFinished) return;

		const lastRound = getCurrentRound(update);
		if (!lastRound) return;

		const winner = findWinner(update, lastRound.track.previewUrl);
		if (!winner) return;

		console.log(lastRound.track, winner, winner.id !== newMe?.id)
		await celebrateTrack(lastRound.track, winner, winner.id !== newMe?.id, true, true);
	}

	const processCelebrations = async (currentGame: GameUpdate | null, update: GameUpdate | null, newMe: Player | null) => {
		if (!update) return;

		if (update.hasFinished) await celebrateEnd(update, newMe);

		const lastSeenRound = getCurrentRound(currentGame)?.number ?? 0;
		const currentRound = getCurrentRound(update)?.number ?? 0;
		if (lastSeenRound >= currentRound) return;

		for (let celebratedRound = lastSeenRound; celebratedRound < currentRound; celebratedRound++) {
			celebrateRound(update, celebratedRound);
		}
	}

	const updateState = ({ game: update, me: newMe }: { game: GameUpdate | null, me: Player | null }) => {
		game = update;
		me = newMe;
		currentRound = getCurrentRound(update);
		currentPlayer = findPlayer(update, currentRound?.playerId);
		isPlaying = !!currentPlayer && currentPlayer.id === me?.id;
	}

	const onUpdate = async ({ game: update, me: newMe }: { game: GameUpdate | null, me: Player | null }) => {
		processCelebrations(game, update, newMe);

		updateState({ game: update, me: newMe });

		if (!update?.hasStarted) return goto(`/game/${gameId}/lobby`);
	}

	onMount(async () => {
		session = new GameSession(gameId);
		session.onUpdate(onUpdate);

		const latestUpdate = session.getLatestUpdate();
		if (latestUpdate) updateState(latestUpdate);
		if (latestUpdate && latestUpdate.game?.rounds.length <= 1) await celebrateStart(latestUpdate.me);
		console.log(latestUpdate)
		if (latestUpdate && latestUpdate.game?.hasFinished) await celebrateEnd(latestUpdate.game, latestUpdate.me);

		await session.initialize();
	});
</script>

<svelte:head>
	<title>Game</title>
	<meta
		name="description"
		content="Shuffle Showdown game. First to collect {game?.songsToWin} songs by placing them in order of release wins!"
	/>
</svelte:head>

<div class="container">
	<div class="overlay" class:hidden={!isCelebrating}>
		<svelte:component this={Celebration} bind:this={celebration} />
	</div>
	<div class="game-interface" class:hidden={isCelebrating}>
		<h1>Round {currentRound?.number}</h1>
		{#if currentPlayer}
			<p>Currently playing: {currentPlayer.id === me?.id ? "you" : currentPlayer.name}</p>
		{/if}
		<Chronology wonTracks={currentPlayer?.wonTracks} onSelect={onAnswerSelect} disabled={!isPlaying}/>
		<AudioPlayer source="{currentRound?.track.previewUrl}" />
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

		h1, p {
			text-align: center;
		}
	}

	.hidden {
		display: none;
	}
</style>

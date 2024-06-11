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

	let celebration: Celebration | null = null;

	const getCurrentRound = (update: GameSessionUpdate | null) : Round | null => {
		if (!update) return null;

		const maxRoundNumber = Math.max(...update.rounds.map((round) => round.number));
		const currentRound = update.rounds.find((round) => round.number === maxRoundNumber);

		return currentRound ?? null;
	}

	const onAnswerSelect = (answer: Answer) => {
		selectedAnswer = answer;
	}

	const onAnswerSubmit = () => {
		if (!selectedAnswer) return;
		isLoading = true;
		session?.submitAnswer(selectedAnswer);
		audioPlayer?.pause();
	}

  const celebrate = ({ game: update, me: newMe }: { game: GameSessionUpdate | null, me: Player | null }) => {
    celebration?.update({
      oldUpdate: gameUpdate,
      newUpdate: update,
      oldMe: me,
      newMe
    });
  };

	const updatePage = async ({ game: update, me: newMe }: { game: GameSessionUpdate | null, me: Player | null }) => {
		gameUpdate = update;
		me = newMe;
		currentRound = getCurrentRound(update);
		currentPlayer = findPlayerInGameSessionUpdate(update, currentRound?.playerId);
		isPlaying = !!currentPlayer && currentPlayer.id === me?.id;
		isLoading = false;

		if (!update?.hasStarted) return goto(`/${gameId}/lobby`);
	}

	onMount(async () => {
		if (!session) session = new GameSession(gameId);
		session.onUpdate((gameUpdate) => {
      celebrate(gameUpdate);
      updatePage(gameUpdate);
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
	<div class="game-interface">
    <div class="game-info">
      <h1>Round {currentRound?.number}</h1>
      {#if currentPlayer}
        <p>Now playing: {currentPlayer.id === me?.id ? "you" : currentPlayer.name} ({currentPlayer.wonTracks?.length}/{gameUpdate?.songsToWin})</p>
      {/if}
    </div>
		<Chronology wonTracks={currentPlayer?.wonTracks} onSelect={onAnswerSelect} disabled={!isPlaying}/>
		<svelte:component this={AudioPlayer} bind:this={audioPlayer} source="{currentRound?.track.previewUrl}" />
		<button class="filled" on:click={onAnswerSubmit} disabled={isLoading}>
			Select answer
		</button>
	</div>
	<svelte:component this={Celebration} bind:this={celebration} />
</div>

<style lang="scss">
	.container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 100%;
		width: 100%;
		position: relative;
	}

	.game-interface {
		flex: 1;
		overflow: hidden;
		display: flex;
		flex-direction: column;
		gap: 1rem;
		align-items: center;
	}

  .game-info {
    h1, p {
      text-align: center;
      margin: 0;
    }
  }
</style>

<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import AudioPlayer from '$lib/components/AudioPlayer.svelte';
	import Chronology from '$lib/components/Chronology.svelte';
	import { GameSession } from '$lib/services/GameSession';
	import type { Player } from '$lib/types/Player';
	import type { Round } from '$lib/types/Round';
	import type { GameUpdate } from '$lib/types/GameUpdate';
	import type { Answer } from '$lib/types/Answer';

	const gameId = $page.params.gameId;
	let session: GameSession | void | null = null;

	let game: GameUpdate | null;
	$: game = null;

	let me: Player | null;
	$: me = null;
	
	let players: Player[];
	$: players = [];

	let currentRound: Round | null;
	$: currentRound = null;

	let selectedAnswer: Answer | null = null;

	const getCurrentRound = (newGame: GameUpdate | null) : Round | null => {
		if (!newGame) return null;

		const maxRoundNumber = Math.max(...newGame.rounds.map((round) => round.number));
		const currentRound = newGame.rounds.find((round) => round.number === maxRoundNumber);

		return currentRound ?? null;
	}

	const onAnswerSelect = (answer: Answer) => {
		selectedAnswer = answer;
	}

	const onAnswerSubmit = () => {
		if (!selectedAnswer) return;
		session?.submitAnswer(selectedAnswer);
	}

	onMount(async () => {
		session = new GameSession(gameId);
		session.onUpdate(({ game: newGame, me: newMe }) => {
			game = newGame;
			me = newMe;
			players = newGame?.players ?? [];
			currentRound = getCurrentRound(newGame);
			console.log('me', me);

			if (!newGame?.isRunning) return goto(`/game/${gameId}/lobby`);
		})
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

<section>
	<h1>Round {currentRound?.number}</h1>
	<Chronology wonTracks={me?.wonTracks} onSelect={onAnswerSelect}/>
	<AudioPlayer source="{currentRound?.track.previewUrl}" />
	<button class="filled" on:click={onAnswerSubmit}>
		Select answer
	</button>
</section>

<style lang="scss">
	section {
		display: grid;
		grid-template-rows: min-content 1fr min-content;
		gap: 1rem;
		justify-content: center;
		align-items: center;

		h1 {
			text-align: center;
		}
	}
</style>

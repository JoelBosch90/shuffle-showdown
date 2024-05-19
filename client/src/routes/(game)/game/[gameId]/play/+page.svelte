<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import AudioPlayer from '$lib/components/AudioPlayer.svelte';
	import { GameSession } from '$lib/services/GameSession';
	import type { Player } from '$lib/types/Player';
	import type { Round } from '$lib/types/Round';
	import type { GameUpdate } from '$lib/types/GameUpdate';

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

	const getCurrentRound = () : Round | null => {
		if (!game) return null;

		const maxRoundNumber = Math.max(...game.rounds.map((round) => round.number));
		const currentRound = game.rounds.find((round) => round.number === maxRoundNumber);

		return currentRound ?? null;
	}

	onMount(async () => {
		session = new GameSession(gameId);
		session.onUpdate(({ game: newGame, me: newMe }) => {
			game = newGame;
			me = newMe;
			players = newGame?.players ?? [];
			currentRound = getCurrentRound();

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
	<AudioPlayer source="{currentRound?.track.previewUrl}" />
</section>

<style lang="scss">
</style>

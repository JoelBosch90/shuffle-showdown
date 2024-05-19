<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { GameSession } from '$lib/services/GameSession';
	import type { Player } from '$lib/types/Player';
	import type { GameUpdate } from '$lib/types/GameUpdate';

	const gameId = $page.params.gameId;
	let session: GameSession | void | null = null;

	let game: GameUpdate | null;
	$: game = null;

	let me: Player | null;
	$: me = null;
	
	let players: Player[];
	$: players = [];

	onMount(async () => {
		session = new GameSession(gameId);
		session.onGameUpdate(({ game: newGame, me: newMe }) => {
			game = newGame;
			me = newMe;
			players = newGame?.players ?? [];

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
	<h1>Game</h1>
</section>

<style lang="scss">
</style>

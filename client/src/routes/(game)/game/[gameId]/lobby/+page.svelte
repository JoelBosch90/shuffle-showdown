<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { API, type Game } from '$lib/services/API';

	const gameId = $page.params.gameId;
	let game: Game | null = null;

	onMount(async () => {
		game = await API.getGame(gameId);

		if (!game) goto(`/game/${gameId}/configure`);
	});

	const players = [];
</script>

<svelte:head>
	<title>Game Lobby</title>
	<meta
		name="description"
		content="Shuffle Showdown game lobby. Share the game code with your friends to join the game."
	/>
</svelte:head>

<section>
	<h1>Game Lobby</h1>
	<h3>Rules:</h3>
	Each team will start with one random song on their timeline, and then take turns playing a random song
	from the playlist. Try to place the song you hear in the correct spot on your timeline. The first team
	to get {game?.songsToWin} songs in the correct spot wins!
	<ul>
		<li>Teams will take turns playing a random song from the playlist.</li>
	</ul>
</section>

<style lang="scss">
</style>

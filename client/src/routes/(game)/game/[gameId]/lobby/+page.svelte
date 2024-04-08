<script>
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { API } from '$lib/services/API';

	const gameId = $page.params.gameId;

	onMount(() => {
		const gameConfig = API.getGame(gameId);

		if (!gameConfig) goto(`/game/${gameId}/configure`);
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
	<ul>
		<li>Teams will take turns playing a random song from the playlist.</li>
	</ul>
</section>

<style lang="scss">
</style>

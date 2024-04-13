<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { API } from '$lib/services/API';
	import { type Game } from '$lib/types/Game';

	const gameId = $page.params.gameId;
	let game: Game | void | null = null;

	const players = [];

	let messages: string[] = [];
	const showMessage = (message: string) => {
		messages = [...messages, message];
	};

	onMount(async () => {
		game = await API.getGame(gameId).catch(() => {
			goto('/game');
		});

		if (!game) goto(`/game/${gameId}/configure`);

		API.SocketConnection.onMessage(({ data }) => showMessage(data));
		API.SocketConnection.start();
	});
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

	<ul>
		{#each messages as message}
			<li>{message}</li>
		{/each}
	</ul>
</section>

<style lang="scss">
</style>

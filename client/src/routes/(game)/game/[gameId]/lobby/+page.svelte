<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { API } from '$lib/services/API';
	import { type Game } from '$lib/types/Game';

	const gameId = $page.params.gameId;
	let url: string | null = null;
	let game: Game | void | null = null;

	const players = [];

	let messages: string[] = [];
	const showMessage = (message: string) => {
		messages = [...messages, message];
	};

	onMount(async () => {
		url = window.location.href;
		game = await API.getGame(gameId).catch(() => {
			goto('/game');
		});

		// TODO: Create player if we don't have one yet.

		if (!game) goto(`/game/${gameId}/configure`);

		API.SocketConnection.onMessage(({ data }) => showMessage(data));
		API.SocketConnection.start(gameId);
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
	Share this link to let your friends join the game: <a href="{url}">{url}</a>
	
	<ul>
		{#each messages as message}
			<li>{message}</li>
		{/each}
	</ul>
</section>

<style lang="scss">
</style>

<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { API } from '$lib/services/API';
	import { type Game } from '$lib/types/Game';
	import { type Player } from '$lib/types/Player';

	const gameId = $page.params.gameId;
	let shareUrl: string | null = null;
	let game: Game | void | null = null;
	let player: Player | void | null = null;

	const players = [];

	let messages: string[] = [];
	const showMessage = (message: string) => {
		messages = [...messages, message];
	};

	onMount(async () => {
		shareUrl = `${window.location.origin}/game/${gameId}/join`;
		game = await API.getGame(gameId).catch(() => {
			return goto('/game');
		});

		if (!game) goto(`/game/${gameId}/configure`);

		player = await API.getPlayer().catch(() => {
			return goto(`/game/${gameId}/join`);
		});

		if (!player) {
			return goto(`/game/${gameId}/join`);
		}

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
	Share this link to let your friends join the game: <a href="{shareUrl}">{shareUrl}</a>
	
	<ul>
		{#each messages as message}
			<li>{message}</li>
		{/each}
	</ul>
</section>

<style lang="scss">
</style>

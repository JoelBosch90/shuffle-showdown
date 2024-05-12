<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { API } from '$lib/services/API';
	import { type Game } from '$lib/types/Game';

	const gameId = $page.params.gameId;
	let game: Game | void | null = null;
	let playerName = '';

	const createPlayer = async () => {
		await API.postPlayer(playerName);

		if (!game?.id) goto(`/game`);

		await goto(`/game/${game?.id}/lobby`);
	};

	onMount(async () => {
		game = await API.getGame(gameId);
		
		if (!game) goto(`/game`);

		// Close any existing socket connection.
		API.closeSocketConnection();

		// Prefill the player's name if possible.
		const player = await API.getPlayer();
		if (player?.name) playerName = player?.name;
	});
</script>

<svelte:head>
	<title>PreGame Lobby</title>
	<meta
		name="description"
		content="Shuffle Showdown pregame lobby. Name yourself to join the lobby."
	/>
</svelte:head>

<section>
	<h1>Join the lobby</h1>

	<p>
		You are joining the lobby to play a game of Shuffle Showdown. Name yourself to join your friends and play!
	</p>

	<ul>
		<li>Playlist name: {game?.playlist.name}</li>
		<li>Playlist owner: {game?.owner.name}</li>
	</ul>

	<label>
		<span>Your name</span><br/>
		<input type="text" placeholder="David Bowie" bind:value={playerName} />
	</label>

	<div class="button-row">
		<button class="filled" on:click={createPlayer}>Join lobby</button>
	</div>
</section>

<style lang="scss">
</style>

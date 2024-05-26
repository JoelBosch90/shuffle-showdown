<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { API } from '$lib/services/API';
	import { type GameUpdate } from '$lib/types/GameUpdate';
	import { type Player } from '$lib/types/Player';

	const gameId = $page.params.gameId;
	let game: GameUpdate | void | null = null;
	let player: Player | void | null = null;
	let playerName = '';
	let isOwner = false;

	const createPlayer = async () => {
		if (!player?.id) await API.postPlayer(playerName);
		else await API.patchPlayer(playerName);

		if (!game?.id) goto(`/game`);

		await goto(`/game/${game?.id}/lobby`);
	};

	onMount(async () => {
		game = await API.getGame(gameId);
		if (!game) goto(`/game`);

		// Prefill the player's name if possible.
		player = await API.getPlayer();
		playerName = player?.name ?? "";
		isOwner = player?.id === game?.owner.id;
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
		{#if isOwner}
			<li>This is your game!</li>
		{:else}
			<li>Playlist owner: {game?.owner.name}</li>
		{/if}
	</ul>

	<form on:submit|preventDefault={createPlayer}>
		<label>
			<span>Your name</span><br/>
			<input type="text" placeholder="David Bowie" bind:value={playerName} />
		</label>

		<div class="button-row">
			<button type="submit" class="filled">Join lobby</button>
		</div>
	</form>
</section>

<style lang="scss">
	form {
		display: flex;
		flex-direction: column;
		align-items: center;
		width: 100%;
		gap: 1rem;

		input {
			box-sizing: border-box;
			width: 100%;
		}
	}
</style>

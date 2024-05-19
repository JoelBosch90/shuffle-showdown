<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { GameSession } from '$lib/services/GameSession';
	import type { Player } from '$lib/types/Player';

	const gameId = $page.params.gameId;
	let shareUrl: string | null = null;
	let session: GameSession | void | null = null;

	let me: Player | null;
	$: me = null;

	let players: Player[];
	$: players = [];

	onMount(async () => {
		shareUrl = `${window.location.origin}/game/${gameId}/join`;
		session = new GameSession(gameId);
		session.onGameUpdate(({ game: newGame, me: newMe }) => {
			me = newMe;
			players = newGame?.players ?? [];

			if (newGame?.isRunning) return goto(`/game/${gameId}/play`);
		})
		await session.initialize();
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
	<h3>Players:</h3>
	<ul class="players">
		{#each players as player}
			<li>
				{#if me?.isOwner && player.id !== me?.id}
				<button on:click={() => session?.kickPlayer(player)}>
					<i class="fa-solid fa-user-slash button kick icon"></i>
				</button>
				{:else}
				<i class="fa-solid fa-user me icon"></i>
				{/if}
				<i class="fa-solid fa-check {player.isConnected ? 'connected' : 'disconnected'} icon"></i>
				<span>
					{player.name}
					{#if player.isOwner}
						<i class="fa-solid fa-crown crown"></i>
					{/if}
				</span>
			</li>
		{/each}
	</ul>

	Share this link to let your friends join the game: <a href="{shareUrl}">{shareUrl}</a>

	{#if me?.isOwner}
		<p>Once everyone is connected, click the button below to start the game.</p>
		<div class="button-row">
			<button class="filled" on:click={session?.startGame}>Start game</button>
		</div>
	{:else}
		<p>Wait for the game owner to start the game.</p>
	{/if}
</section>

<style lang="scss">
	.players {
		display: grid;
		grid-template-columns: min-content min-content max-content;
		row-gap: 0.5rem;
		column-gap: 1rem;
		list-style-type: none;
		padding: 0;

		li {
			display: contents;
		}
	}

	.crown {
		color: var(--yellow);
		margin-left: 1ch;
	}

	.icon { justify-self: center; }
	.connected { color: var(--green); }
	.disconnected { color: var(--red); }
	.kick, .me { color: var(--purple); }
</style>

<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { GameSession } from '$lib/services/GameSession';
	import type { Player } from '$lib/types/Player';
	import LoadingButton from '$lib/components/LoadingButton.svelte';

	const gameId = $page.params.gameId;
	let shareUrl: string | null = null;
	let session: GameSession | void | null = null;
	let isLoading = false;

	let me: Player | null;
	$: me = null;

	let players: Player[];
	$: players = [];

	const onClick = () => {
		isLoading = true;
		session?.startGame();
	}

	onMount(async () => {
		shareUrl = `${window.location.origin}/${gameId}/join`;
		if (!session) session = new GameSession(gameId);
		session.onUpdate(({ game: newGame, me: newMe }) => {
			me = newMe;
			players = newGame?.players ?? [];
			isLoading = false;

			if (newGame?.hasStarted) return goto(`/${gameId}/play`);
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
	{#if players.length > 0}
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
						<span class:anonymous={!player.name}>{player.name || "Anonymous"}</span>
						{#if player.isOwner}
							<i class="fa-solid fa-crown crown"></i>
						{/if}
					</span>
				</li>
			{/each}
		</ul>
	{/if}

	Share this link to let your friends join the game: <a href="{shareUrl}">{shareUrl}</a>

	{#if me?.isOwner}
		<p>Once everyone is connected, click the button below to start the game.</p>
		<div class="button-row">
      <LoadingButton isLoading={isLoading} onClick={onClick}>
        Start game
      </LoadingButton>
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
	.anonymous {
		font-style: italic;
		filter: invert(50%);
	}
</style>

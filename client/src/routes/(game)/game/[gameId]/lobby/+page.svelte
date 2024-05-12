<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { API } from '$lib/services/API';
	import { type Game } from '$lib/types/Game';
	import { type Player } from '$lib/types/Player';
	import { isPlayersUpdateMessage, type PlayersUpdateMessage, type ServerMessage } from '$lib/types/ServerMessage';

	const gameId = $page.params.gameId;
	let shareUrl: string | null = null;
	let game: Game | void | null = null;
	let player: Player | void | null = null;

	let players : Player[];
	$: players = [];

	const handlePlayerUpdate = (message: PlayersUpdateMessage) => {
			players = message.content.map((playerState) => ({
				id: playerState.id,
				name: playerState.name,
				isOwner: playerState.id === game?.owner.id,
				isConnected: playerState.isConnected,
			}));
	};
	const handleMessage = (message: ServerMessage) => {
		if (isPlayersUpdateMessage(message)) return handlePlayerUpdate(message)
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

		API.SocketConnection.onMessage(handleMessage);
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
	<h3>Players:</h3>
	<ul class="players">
		{#each players as player}
			<li>
				<i class="fa-solid fa-check {player.isConnected ? 'connected' : 'disconnected'}"></i>
				{#if player.isOwner}
					<i class="fa-solid fa-crown"></i>
				{/if}
				{player.name}
			</li>
		{/each}
	</ul>

	Share this link to let your friends join the game: <a href="{shareUrl}">{shareUrl}</a>

	<div class="button-row">
		<button class="filled">Start game</button>
	</div>
</section>

<style lang="scss">
	.players {
		list-style-type: none;
		padding: 0;
	}

	i.fa-crown {
		color: var(--yellow);
	}
	i.fa-check.connected {
		color: var(--green);
	}
	i.fa-check.disconnected {
		color: var(--red);
	}
</style>

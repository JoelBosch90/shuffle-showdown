<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { API } from '$lib/services/API';
	import { type Game } from '$lib/types/Game';
	import { type Player } from '$lib/types/Player';
	import { isPlayerKickedMessage, isPlayersUpdateMessage, type PlayersUpdateMessage, type ServerMessage } from '$lib/types/ServerMessage';
	import { ClientMessageType } from '$lib/enums/ClientMessageType';

	const gameId = $page.params.gameId;
	let shareUrl: string | null = null;
	let game: Game | void | null = null;
	let me: Player | void | null = null;

	let players : Player[];
	$: players = [];

	const handlePlayerUpdate = (message: PlayersUpdateMessage) => {
		players = message.payload.map((playerState) => {
			const newState = {
				id: playerState.id,
				name: playerState.name,
				isOwner: playerState.id === game?.owner.id,
				isConnected: playerState.isConnected,
			};

			if (playerState.id === me?.id) {
				me = newState;
			}

			return newState;
		});
	};
	const handleMessage = (message: ServerMessage) => {
		if (isPlayersUpdateMessage(message)) return handlePlayerUpdate(message);
		if (isPlayerKickedMessage(message)) return goto('/game');

	};
	const kickPlayer = (playerToKick: Player) => {
		API.sendSocketMessage({
			type: ClientMessageType.KickPlayer,
			payload: playerToKick.id,
		});
	};

	onMount(async () => {
		shareUrl = `${window.location.origin}/game/${gameId}/join`;
		game = await API.getGame(gameId).catch(() => {
			return goto('/game');
		});

		if (!game) goto(`/game/${gameId}/configure`);

		me = await API.getPlayer().catch(() => {
			return goto(`/game/${gameId}/join`);
		});

		if (!me) {
			return goto(`/game/${gameId}/join`);
		}

		API.onSocketMessage(handleMessage);
		API.startSocketConnection(gameId);
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
				<button on:click={() => kickPlayer(player)}>
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

	<div class="button-row">
		<button class="filled">Start game</button>
	</div>
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

<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { API, type Game } from '$lib/services/API';

	const gameId = $page.params.gameId;
	let game: Game | null = null;

	const players = [];

  let connected = false;
  let connection: WebSocket | undefined = undefined;

  let messages: string[] = [];
  const showMessage = (message: string) => {
    messages = [ ...messages, message ];
  };

  const connectWebSocket = () => {
    if (connected) return;

    const protocol = window.location.protocol.endsWith('s:') ? 'wss:' : 'ws:';
    connection = new WebSocket(`${protocol}//${window.location.host}/api/v1/ws`);

    connection.addEventListener('open', (event) : void => {
      connected = true;
      showMessage('Connected.')
    });

    connection.addEventListener('close', (event) : void => {
      connected = false;
      showMessage('Disconnected.')
    });

    connection.addEventListener('message', (event) : void => {
      showMessage(`Message received: ${event.data}`);
    });
  }

	onMount(async () => {
		game = await API.getGame(gameId);

		if (!game) goto(`/game/${gameId}/configure`);

    connectWebSocket();
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

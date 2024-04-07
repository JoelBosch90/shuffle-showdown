<script>
	import { goto } from '$app/navigation';

	let playlist = '';

	const createGame = async () => {
		const playlistId = parseInt(playlist);

		const response = await fetch('/api/v1/game', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ playlist: playlistId })
		});

		if (!response.ok) return;

		const { data } = await response.json();

		await goto(`/game/${data.id}/configure`);
	};
</script>

<svelte:head>
	<title>Select Playlist</title>
	<meta
		name="description"
		content="Select a Spotify Playlist to create a new Shuffle Showdown game."
	/>
</svelte:head>

<section>
	<h1>Select your playlist</h1>

	<input placeholder="Paste your Spotify Playlist link here..." bind:value={playlist} />

	<button class="filled" on:click={createGame}>Select playlist</button>
</section>

<style lang="scss">
	input {
		width: 100%;
	}
</style>

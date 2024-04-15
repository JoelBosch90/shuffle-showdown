<script>
	import { goto } from '$app/navigation';
	import { API } from '$lib/services/API';

	let playListId = '';

	const createGame = async () => {
		const { data } = await API.postGame(playListId.split('/playlist/')[1].split('?')[0]);

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

	<input placeholder="Paste your Spotify Playlist link here..." bind:value={playListId} />

	<button class="filled" on:click={createGame}>Select playlist</button>
</section>

<style lang="scss">
	input {
		box-sizing: border-box;
		width: 100%;
	}
</style>

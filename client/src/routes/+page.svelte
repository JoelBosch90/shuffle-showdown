<script>
	import { goto } from '$app/navigation';
	import LoadingButton from '$lib/components/LoadingButton.svelte';
	import { API } from '$lib/services/API';

	let playList = '';
	let isLoading = false;

	const createGame = async () => {
		isLoading = true;
		const game = await API.postGame(playList);
		isLoading = false;
		if (!game) return;

		await goto(`/${game.id}/join`);
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
	<form on:submit|preventDefault={createGame}>
		<h1>Select your playlist</h1>

		<input placeholder="Paste your Spotify Playlist link here..." bind:value={playList} />

    <LoadingButton type="submit" isLoading={isLoading}>
      Select playlist
    </LoadingButton>
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

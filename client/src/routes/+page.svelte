<script>
	import { goto } from '$app/navigation';
	import { API } from '$lib/services/API';
  import { showToast } from '$lib/store/toasts';
  import { ToastType } from '$lib/types/Toast';

	let playList = '';

	const createGame = async () => {
		const game = await API.postGame(playList);

		await goto(`/${game.id}/join`);
	};

  const createToast = () => {
    showToast({
      type: ToastType.Success,
      message: 'Happy stuff',
    })
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

		<button type="submit" class="filled">Select playlist</button>
    <button class="filled" on:click={createToast} type="button">Show toast</button>
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

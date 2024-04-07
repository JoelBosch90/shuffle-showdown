<script>
	import { goto } from '$app/navigation';

	const createGame = async () => {
		const textarea = document.querySelector('textarea');
		const playlist = parseInt(textarea?.value ?? '');

		const response = await fetch('/api/v1/game', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ playlist })
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

	<textarea cols="1" placeholder="Paste your Spotify Playlist link here..." />

	<button on:click={createGame}>Select playlist</button>
</section>

<style lang="scss">
	section {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 1rem;

		textarea {
			box-sizing: border-box;
			width: 100%;
			resize: none;
			padding: 0.5rem;
		}

		button {
			padding: 0.5rem 1rem;
			color: var(--white);
			background-color: var(--purple);
			border-radius: var(--border-radius);
		}
	}
</style>

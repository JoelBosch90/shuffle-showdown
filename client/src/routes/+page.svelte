<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import LoadingButton from '$lib/components/LoadingButton.svelte';
	import { API } from '$lib/services/API';
	import type { MostPlayedPlaylist } from '$lib/types/MostPlayedPlaylist';

	let playlistInput = '';
	let isLoading = false;
	let mostPlayedPlaylists: MostPlayedPlaylist[] = [];

	const createGame = async (playlist: string) => {
		isLoading = true;
		const game = await API.postGame(playlist);
		isLoading = false;
		if (!game) return;

		await goto(`/${game.id}/join`);
	};

	onMount(async () => {
		mostPlayedPlaylists = (await API.getMostPlayedPlaylists()) ?? [];
	});
</script>

<svelte:head>
	<title>Select Playlist</title>
	<meta
		name="description"
		content="Select a Spotify Playlist to create a new Shuffle Showdown game."
	/>
</svelte:head>

<section>
	<div class="most-played">
		<h1>Most played</h1>
		<ol>
			{#each mostPlayedPlaylists as playlist}
				{#if playlist.playlist.name}
					<li>
						<button on:click={() => createGame(playlist.playlist.id)}>
							{playlist.playlist.name} ({playlist.gamesPlayed} games)
						</button>
					</li>
				{/if}
			{/each}
		</ol>
	</div>

	<form on:submit|preventDefault={() => createGame(playlistInput)}>
		<h1>Select a new playlist</h1>

		<input
			type="text"
			placeholder="Paste your Spotify Playlist link here..."
			bind:value={playlistInput}
		/>

		<LoadingButton type="submit" {isLoading}>Select playlist</LoadingButton>
	</form>
</section>

<style lang="scss">
	section {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 2rem;
	}

	.most-played {
		display: flex;
		flex-direction: column;
		align-items: center;

		h1 {
			margin: 0;
		}

		li {
			padding: 0;

			button {
				margin: 0.25rem 0;
				color: var(--purple);
			}
		}
	}

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

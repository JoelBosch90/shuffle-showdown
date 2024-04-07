<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { gameConfigStore, type GameConfig } from '$lib/stores/gameConfigStore';

	const settings: GameConfig = {
		gameId: $page.params.gameId,
		songsToWin: 10,
		titleRequired: false,
		artistRequired: false,
		releaseYearRequired: true
	};

	const configureGame = async () => {
		gameConfigStore.set(settings);

		await goto(`/game/${settings.gameId}/start`);
	};
</script>

<svelte:head>
	<title>Configure Game</title>
	<meta name="description" content="Configure your Shuffle Showdown game." />
</svelte:head>

<section>
	<h1>Configure your game</h1>

	<input type="number" placeholder="Number of songs to win" bind:value={settings.songsToWin} />

	<label>
		<input type="checkbox" bind:checked={settings.titleRequired} />
		Title required
	</label>

	<label>
		<input type="checkbox" bind:checked={settings.artistRequired} />
		Artist required
	</label>

	<label>
		<input type="checkbox" bind:checked={settings.releaseYearRequired} />
		Release year required
	</label>

	<button on:click={configureGame}>Create game</button>
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

<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { type GameConfig } from '$lib/types/GameConfig';
	import { API } from '$lib/services/API';

	const settings: GameConfig = {
		gameId: $page.params.gameId,
		songsToWin: 10,
		titleRequired: false,
		artistRequired: false,
		releaseYearRequired: true
	};

	const configureGame = async () => {
		await API.patchGame(settings);

		await goto(`/game/${settings.gameId}/lobby`);
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

	<button class="filled" on:click={configureGame}>Create game</button>
</section>

<style lang="scss">
</style>

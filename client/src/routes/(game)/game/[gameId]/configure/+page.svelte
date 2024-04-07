<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { gameConfigStore, type GameConfig } from '$lib/stores/gameConfigStore';

	enum PredefinedGameMode {
		normal = 'normal',
		hard = 'hard'
	}
	type GameSettings = Omit<GameConfig, 'gameId'>;

	const predefinedSettings: Record<PredefinedGameMode, GameSettings> = {
		normal: {
			songsToWin: 10,
			titleRequired: false,
			artistRequired: false,
			releaseYearRequired: true
		},
		hard: {
			songsToWin: 10,
			titleRequired: true,
			artistRequired: true,
			releaseYearRequired: true
		}
	};

	const gotoCustomGame = async () => {
		await goto(`/game/${$page.params.gameId}/configure/custom`);
	};

	const configureGame = async (settings: GameSettings) => {
		const gameId = $page.params.gameId;

		const config = { ...settings, gameId };
		gameConfigStore.set(config);

		await goto(`/game/${gameId}/start`);
	};
</script>

<svelte:head>
	<title>Configure Game</title>
	<meta name="description" content="Configure your Shuffle Showdown game." />
</svelte:head>

<section>
	<h1>Configure your game</h1>

	<button class="filled" on:click={() => configureGame(predefinedSettings[PredefinedGameMode.hard])}
		>Hard mode</button
	>

	<button
		class="filled"
		on:click={() => configureGame(predefinedSettings[PredefinedGameMode.normal])}
		>Normal mode</button
	>

	<button class="filled" on:click={gotoCustomGame}>Custom mode</button>
</section>

<style lang="scss">
</style>

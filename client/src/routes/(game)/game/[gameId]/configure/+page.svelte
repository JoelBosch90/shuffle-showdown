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

	<button on:click={() => configureGame(predefinedSettings[PredefinedGameMode.hard])}
		>Hard mode</button
	>

	<button on:click={() => configureGame(predefinedSettings[PredefinedGameMode.normal])}
		>Normal mode</button
	>

	<button on:click={gotoCustomGame}>Custom mode</button>
</section>

<style lang="scss">
	section {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 1rem;

		button {
			padding: 0.5rem 1rem;
			color: var(--white);
			background-color: var(--purple);
			border-radius: var(--border-radius);
		}
	}
</style>

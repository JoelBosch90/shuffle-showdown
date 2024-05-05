<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { type GameConfig } from '$lib/types/GameConfig';
	import { API } from '$lib/services/API';

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

		await API.patchGame({ ...settings, gameId });

		await goto(`/game/${gameId}/lobby`);
	};
</script>

<svelte:head>
	<title>Configure Game</title>
	<meta name="description" content="Configure your Shuffle Showdown game." />
</svelte:head>

<section>
	<h1>Configure your game</h1>

	<label>
		<span>Your name</span><br/>
		<input type="text" placeholder="David Bowie" />
	</label>


	<h3>Start the game in:</h3>

	<div class="button-row">
		<button class="filled" on:click={() => configureGame(predefinedSettings[PredefinedGameMode.hard])}>
			Hard mode
		</button>
	
		<button class="filled" on:click={() => configureGame(predefinedSettings[PredefinedGameMode.normal])}>
			Normal mode
		</button>
	
		<button class="filled" on:click={gotoCustomGame}>Custom mode</button>
	</div>
</section>

<style lang="scss">
	label, input {
		box-sizing: border-box;
		width: 100%;
	}

	label {
		color: var(--purple);
		font-weight: bold;

		span {
			margin-left: 1ch;
		}
	}

	.button-row {
		display: flex;
		justify-content: space-between;
		gap: 1rem;
	}
</style>

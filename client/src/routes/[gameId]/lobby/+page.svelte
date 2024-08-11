<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { GameSession } from '$lib/services/GameSession';
	import type { Player } from '$lib/types/Player';
	import type { PlaylistLoadingUpdate } from '$lib/types/PlaylistLoadingUpdate';
	import LoadingButton from '$lib/components/LoadingButton.svelte';
	import ProgressBar from '$lib/components/ProgressBar.svelte';
	import QRCode from 'qrcode';
	import { showToast } from '$lib/store/toasts';
	import { ToastType } from '$lib/enums/ToastType';

	const gameId = $page.params.gameId;

	let canvas: HTMLCanvasElement | null = null;
	let url: string | null = null;
	let session: GameSession | void | null = null;
	let isLoading = false;
	let isDisabled = true;
	let latestPlaylistLoadingUpdate: PlaylistLoadingUpdate | null = null;
	let tracksTotal = 0;
	let tracksLoaded = 0;
	let tracksChecked = 0;
	let doneLoading = false;
	let doneChecking = false;
	$: tracksTotal = latestPlaylistLoadingUpdate?.tracksTotal ?? 0;
	$: tracksLoaded = latestPlaylistLoadingUpdate?.tracksLoaded ?? 0;
	$: tracksChecked = latestPlaylistLoadingUpdate?.tracksChecked ?? 0;
	$: tracksLoadingLabel = `Loaded ${tracksLoaded} out of ${tracksTotal} tracks`;
	$: tracksCheckingLabel = `Checked release dates for ${tracksChecked} out of ${tracksLoaded} tracks`;
	$: doneLoading = latestPlaylistLoadingUpdate?.finishedLoading ?? false;
	$: doneChecking = latestPlaylistLoadingUpdate?.finishedChecking ?? false;

	let owner: Player | null;
	$: owner = players.find((player) => player.isOwner) ?? null;

	let me: Player | null;
	$: me = null;

	let players: Player[];
	$: players = [];

	let playlistName = '';

	const onClick = () => {
		isLoading = true;
		session?.startGame();
	};

	const canShare = () => typeof navigator.share === 'function';

	const shareUrl = () =>
		navigator.share({
			title: 'Shuffle Showdown',
			text: `Join ${owner?.name ? `${owner.name}'s` : 'our'} game of Shuffle Showdown${playlistName ? ` with the '${playlistName}' playlist` : ''}!`,
			url: url ?? undefined
		});

	const copyUrl = () => {
		if (!url) return;
		navigator.clipboard.writeText(url);
	};

	onMount(async () => {
		url = `${window.location.origin}/${gameId}/join`;
		if (!session) session = new GameSession(gameId);

		QRCode.toCanvas(canvas, url, {
			scale: 1,
			width: 1280,
			margin: 0,
			color: {
				dark: '#642CA9',
				light: '#F0F2EF'
			}
		});

		session.onUpdate(({ game: newGame, me: newMe }) => {
			me = newMe;
			players = newGame?.players ?? [];
			isLoading = false;
			isDisabled = !(newGame?.isReadyToStart ?? false);
			playlistName = newGame?.playlist?.name ?? '';

			if (newGame?.hasStarted) return goto(`/${gameId}/play`);
		});
		session.onPlaylistLoadingUpdate(({ playlistLoadingUpdate }) => {
			latestPlaylistLoadingUpdate = playlistLoadingUpdate;
		});

		await session.initialize();
	});
</script>

<svelte:head>
	<title>Game Lobby</title>
	<meta
		name="description"
		content="Shuffle Showdown game lobby. Share the game code with your friends to join the game."
	/>
</svelte:head>

<section>
	<div class="title">
		<h1>Game Lobby</h1>
		{#if playlistName}
			<p>{playlistName}</p>
		{/if}
	</div>

	<div class="content">
		<div class="share-options">
			<h3>
				Invite link
				{#if canShare()}
					<button class="share" on:click={shareUrl}>
						<i class="fa-solid fa-share-nodes icon"></i>
					</button>
				{/if}
			</h3>

			<canvas bind:this={canvas}></canvas>

			<button class="copy" on:click={copyUrl}>
				<span class="share-url">{url}</span>
				<i class="fa-solid fa-copy icon"></i>
			</button>
		</div>

		<div class="game-options">
			{#if players.length > 0}
				<ul class="players">
					{#each players as player}
						<li style="--connection-color: var(--{player.isConnected ? 'green' : 'red'});">
							{#if me?.isOwner && player.id !== me?.id}
								<button on:click={() => session?.kickPlayer(player)}>
									<i class="fa-solid fa-ban kick icon"></i>
								</button>
							{:else if player?.isOwner}
								<i class="fa-solid fa-crown crown icon"></i>
							{:else}
								<i class="fa-solid fa-user me icon"></i>
							{/if}
							<span>
								<span class:anonymous={!player.name}>{player.name || 'Anonymous'}</span>
							</span>
						</li>
					{/each}
				</ul>
			{/if}

			<div class="button-row">
				{#if me?.isOwner}
					<LoadingButton {isLoading} {isDisabled} {onClick}>Start game</LoadingButton>
				{:else}
					<p>Wait for {owner?.name ?? 'the owner'} to start the game.</p>
				{/if}
			</div>
		</div>

		<div class="progress-bars">
			<ProgressBar
				bind:total={tracksTotal}
				bind:current={tracksLoaded}
				bind:label={tracksLoadingLabel}
				bind:done={doneLoading}
			/>
			<ProgressBar
				bind:total={tracksLoaded}
				bind:current={tracksChecked}
				bind:label={tracksCheckingLabel}
				bind:done={doneChecking}
			/>
		</div>
	</div>
</section>

<style lang="scss">
	section {
		--gap: clamp(0.5rem, 3cqh, 2rem);
		--title-height: 3rem;

		container-name: section;
		container-type: inline-size;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: var(--gap);

		box-sizing: border-box;
		max-height: 100%;
		overflow: hidden;
	}

	.title {
		display: flex;
		flex-shrink: 1;
		flex-direction: column;
		align-items: center;
		height: var(--title-height);

		h1,
		p {
			margin: 0;
		}
	}

	.content {
		display: grid;
		grid-template-areas:
			'share-options'
			'game-options'
			'progress-bars';
		justify-items: center;
		gap: var(--gap);

		box-sizing: border-box;
		max-height: 100%;
		width: 100%;
		overflow: hidden;

		.share-options,
		.game-options {
			display: flex;
			flex-direction: column;
			align-items: center;

			max-height: 100%;
			width: 100%;
			overflow: hidden;
		}

		.share-options {
			--share-title-height: 1.5rem;
			--share-url-height: 1.25rem;
			--min-qr-code-size: 120px;
			--max-qr-code-size: 400px;
			$vertical-text-height: calc(
				3 * var(--gap) + var(--title-height) + var(--share-title-height) + var(--share-url-height)
			);
			$qr-code-size: clamp(
				var(--min-qr-code-size),
				min(80cqw, min(50cqw, calc(50cqh - $vertical-text-height))),
				var(--max-qr-code-size)
			);

			grid-area: share-options;
			container-name: share-options;
			min-height: min-content;
			min-width: min-content;
			width: $qr-code-size;

			h3 {
				margin: 0;
				height: var(--share-title-height);
				display: flex;
				flex-direction: row;
			}

			canvas,
			button {
				margin: 0.5rem;
			}

			canvas {
				max-width: $qr-code-size;
				max-height: $qr-code-size;
			}

			button {
				display: flex;
				flex-direction: row;
				align-items: flex-start;
				line-height: 1.25;
				margin-top: 0;

				.share-url {
					display: -webkit-box;
					font-size: 0.75rem;
					line-clamp: 2;
					-webkit-line-clamp: 2;
					-webkit-box-orient: vertical;
					word-break: break-all;
					overflow: hidden;
					text-overflow: ellipsis;

					margin-right: 2ch;
				}
			}

			.copy,
			.share {
				i {
					color: var(--purple);
				}
			}
		}

		.game-options {
			grid-area: game-options;
			gap: 1rem;
			justify-content: space-around;

			.players {
				display: grid;
				grid-template-columns: min-content max-content;
				row-gap: 0.5rem;
				column-gap: 1rem;
				list-style-type: none;
				padding: 0;
				margin: 0;
				width: calc(100% - 2 * var(--gap));
				height: max-content;
				justify-content: center;
				overflow-y: auto;

				li {
					--connection-color: var(--green);

					display: contents;

					.crown {
						color: var(--yellow);
					}

					.kick,
					.me {
						color: var(--purple);
					}

					.crown,
					.kick,
					.me {
						position: relative;
						min-width: calc(1em + 1ch);
						min-height: 1em;

						&::before {
							position: absolute;
							top: 50%;
							left: 50%;
							transform: translate(-50%, -50%);
						}

						&::after {
							--connection-dot-side: 1ch;

							position: absolute;
							bottom: 0;
							right: 0;

							content: '';
							width: var(--connection-dot-side);
							height: var(--connection-dot-side);

							border-radius: 50%;
							background-color: var(--connection-color);
						}
					}
				}

				.anonymous {
					font-style: italic;
					filter: invert(50%);
				}
			}
		}

		.progress-bars {
			grid-area: progress-bars;
			display: flex;
			flex-direction: column;
			gap: 0.25rem;
			width: 100%;
			box-sizing: border-box;
		}
	}

	.icon {
		justify-self: center;
	}

	@container section (min-width: 450px) {
		.content {
			--gap: clamp(0.5rem, 3cqw, 2rem);
			--padding: clamp(0.5rem, 9cqh, 6rem);
			width: 100%;
			box-sizing: border-box;

			display: grid;
			grid-template-areas:
				'share-options game-options'
				'progress-bars progress-bars';
			grid-template-columns: 1fr 1fr;

			.share-options {
				$vertical-text-height: calc(
					3 * var(--gap) + var(--title-height) + var(--share-title-height) + var(--share-url-height)
				);
				$qr-code-size: clamp(
					var(--min-qr-code-size),
					min(80cqw, min(50cqmin, calc(100cqh - $vertical-text-height))),
					var(--max-qr-code-size)
				);

				canvas {
					max-width: $qr-code-size;
					max-height: $qr-code-size;
				}
			}

			.game-options {
				box-sizing: border-box;
				height: 100%;
				padding: var(--padding) 0;
			}
		}
	}
</style>

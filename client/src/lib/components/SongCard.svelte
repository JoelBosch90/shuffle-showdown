<script lang="ts">
	import type { Artist } from '$lib/types/Artist';

	export let cardColor: string = '';
	export let cardIcon: string = '';
	export let disabled: boolean = false;
	export let normalizedIndex: number = 0;
	export let releaseYear: string = '???';
	export let trackName: string = '';
	export let artists: Artist[] = [];
	export let isGuess: boolean = false;
	export let isRevealing: boolean = false;
	export let isWon: boolean = false;

	const joinArtists = (artists: Artist[] = []) => {
		return artists.map((artist) => artist.name).join(', ');
	};

	$: artistString = joinArtists(artists);
</script>

<li
	class:disabled
	class="{isGuess ? 'guess' : ''} {isRevealing ? 'revealing' : ''} {isRevealing && isWon
		? 'correct'
		: ''} {isRevealing && !isWon ? 'wrong' : ''}"
	style="--normalized-index: {normalizedIndex}; --card-color: var({cardColor}, --gray-dark);"
>
	{#if cardIcon}
		<i class={`fa-solid fa-${cardIcon} player-icon`}></i>
	{/if}

	<h2 title={releaseYear}>{releaseYear}</h2>

	{#if trackName}
		<p class="track" title={trackName}>{trackName}</p>
	{/if}

	{#if artistString}
		<p class="artist" title={artistString}>{artistString}</p>
	{/if}
</li>

<style lang="scss">
	li {
		--normalized-index: 0;
		--card-padding: 1rem;
		--card-font-size: 12cqw;
		--card-min-height: 6rem;
		--card-scale-height: 50cqmin;
		--card-max-height: 16rem;
		--card-aspect-ratio: 1 / 1.25;
		--card-box-shadow-space: 0.75rem;
		--card-color: var(--gray-dark);
		--center-card-margin-percentage: 0.33;

		$normalized-index: var(--normalized-index);
		$distance-from-guess-card: max($normalized-index, -1 * $normalized-index);
		$card-height: clamp(var(--card-min-height), var(--card-scale-height), var(--card-max-height));
		$card-width: calc($card-height * var(--card-aspect-ratio));

		display: flex;
		box-sizing: border-box;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		gap: 0.2rem;
		opacity: 1;
		list-style: none;
		container-type: size;
		container-name: li;

		position: absolute;
		top: 0;
		left: 0;
		aspect-ratio: var(--card-aspect-ratio);
		height: $card-height;
		z-index: calc(var(--card-level) - $distance-from-guess-card);

		padding: var(--card-padding);
		border-radius: var(--card-padding);
		box-shadow: 0 0 var(--card-box-shadow-space) rgba(0, 0, 0, 0.25);
		color: var(--white);
		background-color: var(--card-color, --white);
		overflow: hidden;

		transition:
			transform var(--animation-speed-quick) ease-in,
			z-index var(--animation-speed-quick) ease-in,
			opacity var(--animation-speed-quick) ease-in;

		$direction: clamp(-1, $normalized-index, 1);
		$non-zero-index: calc($distance-from-guess-card + var(--center-card-margin-percentage));
		$fraction: calc(1 / $non-zero-index);
		$clamped-fraction: clamp(0, $fraction, 1);
		$scale: clamp(0, $clamped-fraction, 1);
		$index-offset: calc($scale - 1);
		$card-offset: calc($index-offset * 50 * $direction);
		$size-scale: calc($scale * 0.5 + 0.5);

		@function calculated-offset($side-unit, $card-side-length) {
			@return calc(
				$card-offset * $side-unit - $direction * $card-side-length *
					var(--center-card-margin-percentage)
			);
		}
		$card-vertical-offset: calculated-offset(1cqh, $card-height);
		$card-horizontal-offset: calculated-offset(1cqw, $card-width);

		@function scaled-value($value) {
			@return calc($value * $size-scale);
		}
		$card-scaled-height: scaled-value($card-height);
		$card-scaled-width: scaled-value($card-width);
		$card-scaled-box-shadow-space: scaled-value(var(--card-box-shadow-space));

		@function reverse-scaled-value($value) {
			@return calc($value * (1 - $size-scale));
		}
		$card-scaling-height-loss: reverse-scaled-value($card-height);
		$card-scaling-width-loss: reverse-scaled-value($card-width);

		@function start-value($card-side-length) {
			@return calc(
				0% - reverse-scaled-value($card-side-length) * 0.5 + $card-scaled-box-shadow-space
			);
		}
		$card-at-top: start-value($card-height);
		$card-at-left: start-value($card-width);

		@function middle-value($side-unit, $card-side-length) {
			@return calc(
				50 * $side-unit -
					(scaled-value($card-side-length) + reverse-scaled-value($card-side-length)) * 0.5
			);
		}
		$card-at-vertical-center: middle-value(1cqh, $card-height);
		$card-at-horizontal-center: middle-value(1cqw, $card-width);

		@function end-value($side-unit, $card-side-length) {
			@return calc(
				100 * $side-unit -
					(
						$card-side-length - reverse-scaled-value($card-side-length) * 0.5 +
							$card-scaled-box-shadow-space
					)
			);
		}
		$card-at-bottom: end-value(1cqh, $card-height);
		$card-at-right: end-value(1cqw, $card-width);

		$vertical-transform: translateY(
			clamp($card-at-top, calc($card-at-vertical-center + $card-vertical-offset), $card-at-bottom)
		);
		$horizontal-transform: translateX(
			clamp(
				$card-at-left,
				calc($card-at-horizontal-center - $card-horizontal-offset),
				$card-at-right
			)
		);
		transform: $vertical-transform $horizontal-transform scale($size-scale);

		p,
		h2 {
			max-width: 100%;
			margin: 0;
			overflow: hidden;
			text-overflow: ellipsis;
			text-align: center;
		}

		p {
			display: -webkit-box;
			line-clamp: 2;
			-webkit-line-clamp: 2;
			-webkit-box-orient: vertical;
			font-size: var(--card-font-size);
		}

		h2 {
			font-size: calc(var(--card-font-size) * 3.2);
			padding: 0.1em;
		}

		&.disabled {
			color: var(--gray-dark);
		}

		&.guess {
			opacity: 50%;
		}

		.player-icon {
			position: absolute;
			top: var(--card-padding);
			font-size: var(--card-font-size);
		}

		&.revealing {
			transition:
				background-color var(--animation-speed-quick) ease-out var(--animation-speed-slow),
				transform var(--animation-speed-slow) cubic-bezier(0, 0.25, 1, 0.25),
				opacity var(--animation-speed-quick) ease-out;

			transform: $vertical-transform $horizontal-transform scale(calc($size-scale + 0.5));

			&.wrong {
				background-color: var(--red);
			}
		}
	}
</style>

<script lang="ts">
  import type { WonTrack } from '$lib/types/WonTrack';
  import type { Artist } from '$lib/types/Artist';
  import type { Answer } from '$lib/types/Answer';
	import { onMount } from 'svelte';

  export let disabled: boolean = false;
  export let wonTracks: WonTrack[] = [];
  export let onSelect: (answer: Answer) => void;

  interface Card {
    releaseYear: string;
    name?: string;
    artists?: string;
    isGuess?: boolean;
  }

  let container: HTMLOListElement;

  const UNKNOWN_RELEASE_YEAR = "???";
  const guessCard: Card = {
    releaseYear: UNKNOWN_RELEASE_YEAR,
    isGuess: true,
  }

  const joinArtists = (artists: Artist[] = []) => {
    return artists.map((artist) => artist.name).join(', ');
  };

  const selectAnswer = () => {
    if (disabled) return;

    const cardBefore = guessIndex > 0 ? cards[guessIndex - 1] : undefined;
    const cardAfter = guessIndex < cards.length - 1 ? cards[guessIndex + 1] : undefined;

    onSelect({
      afterReleaseYear: cardBefore ? parseInt(cardBefore?.releaseYear ?? "") : undefined,
      beforeReleaseYear: cardAfter ? parseInt(cardAfter?.releaseYear ?? "") : undefined,
    });
  };

  const trackToCard = ({ track }: WonTrack) : Card => ({
    releaseYear: track.releaseYear?.toString() ?? UNKNOWN_RELEASE_YEAR,
    name: track.name,
    artists: joinArtists(track.artists),
  });

  const sortWonTracks = (a: WonTrack, b: WonTrack) => {
    const aReleaseYear = a.track.releaseYear ?? 0;
    const bReleaseYear = b.track.releaseYear ?? 0;
    return aReleaseYear - bReleaseYear;
  };

  const onWheelEvent = (event: WheelEvent) => {
    const direction = (event.deltaY + event.deltaX) > 0 ? 1 : -1;

    guessIndex = clampGuessIndex(guessIndex + direction);
  };

  const clampGuessIndex = (newGuessIndex: number) => {
    return Math.min(Math.max(0, newGuessIndex), trackCards.length);
  };

  const getClientLocation = (event: MouseEvent | TouchEvent) => {
    if (event instanceof MouseEvent) return { x: event.clientX, y: event.clientY };
    return { x: event.touches[0].clientX, y: event.touches[0].clientY };
  };

  const onMoveStart = (event: MouseEvent | TouchEvent) => {
    const { x: startX, y: startY } = getClientLocation(event);
    const startGuessIndex = guessIndex;

    const onMove = (event: MouseEvent | TouchEvent) => {
      const { x: currentX, y: currentY } = getClientLocation(event);
      const percentageOfScreenToTravelToMoveThroughAllOptions = 0.8;
      const heightPercentageTraveled = (currentX - startX) / (container.clientWidth * percentageOfScreenToTravelToMoveThroughAllOptions);
      const widthPercentageTraveled = (currentY - startY) / (container.clientHeight * percentageOfScreenToTravelToMoveThroughAllOptions);

      const distance = Math.hypot(widthPercentageTraveled, heightPercentageTraveled);
      const indexChange = Math.floor(distance * trackCards.length);
      const direction = (widthPercentageTraveled - heightPercentageTraveled) > 0 ? 1 : -1;

      guessIndex = clampGuessIndex(startGuessIndex + direction * indexChange);
    };

    const onMoveEnd = () => {
      window.removeEventListener('touchmove', onMove);
      window.removeEventListener('touchend', onMoveEnd);
      window.removeEventListener('touchcancel', onMoveEnd);

      window.removeEventListener('mousemove', onMove);
      window.removeEventListener('mouseup', onMoveEnd);
    };

    window.addEventListener('touchmove', onMove);
    window.addEventListener('touchend', onMoveEnd);
    window.addEventListener('touchcancel', onMoveEnd);

    window.addEventListener('mousemove', onMove);
    window.addEventListener('mouseup', onMoveEnd);
  };

  let trackCards: Card[];
  $: trackCards = wonTracks.toSorted(sortWonTracks).map(trackToCard);

  let guessIndex: number;
  $: guessIndex = Math.ceil(trackCards.length / 2);

  let cards: Card[];
  $: cards = [...trackCards.slice(0, guessIndex), guessCard, ...trackCards.slice(guessIndex)];
  $: cards, selectAnswer();

  onMount(() => {
    container.addEventListener('wheel', onWheelEvent);
    container.addEventListener('mousedown', onMoveStart);
    container.addEventListener('touchstart', onMoveStart);
  });
</script>

<ol class="chronology" class:disabled={disabled} bind:this="{container}">
  {#each cards as card, cardIndex}
    <li class="card {card.isGuess ? 'guess' : ''}" style="--normalized-index: {cardIndex - guessIndex}">
      <h2>{card.releaseYear}</h2>

      {#if card.name}
        <p class="track">{card.name}</p>
      {/if}

      {#if card.artists}
        <p class="artist">{card.artists}</p>
      {/if}
    </li>
  {/each}
</ol>

<style lang="scss">
  .chronology {
    container-type: size;
    container-name: chronology;
    display: flex;
    flex-grow: 1;
    width: 100%;
    margin: 0;
    padding: 0;
    overflow: hidden;
    position: relative;
    box-sizing: border-box;

    cursor: grab;
    user-select: none;
    -webkit-user-select: none;

    &.disabled {
      .card {
        color: var(--gray-dark);
      }
    }

    .card {
      --normalized-index: 0;
      --card-padding: 1rem;
      --card-min-height: 6rem;
      --card-scale-height: 35cqh;
      --card-max-height: 16rem;
      --card-font-size: 12cqw;
      --card-aspect-ratio: 1 / 1.25;
      --card-box-shadow-space: 0.75rem;
      --card-box-shadow-margin: 5px;
      --default-vertical-distance: 15cqh;
      --default-horizontal-distance: 15cqw;
      --distance-step: 0.125;

      $distance-from-guess-card: max(var(--normalized-index), -1 * var(--normalized-index));
      $card-height: clamp(var(--card-min-height), var(--card-scale-height), var(--card-max-height));
      $card-width: calc($card-height * var(--card-aspect-ratio));

      display: flex;
      box-sizing: border-box;
      container-type: size;
      container-name: card;

      position: absolute;
      top: 0;
      left: 0;

      aspect-ratio: var(--card-aspect-ratio);
      height: $card-height;
      z-index: calc(var(--card-level) - $distance-from-guess-card);

      $direction: clamp(-1, var(--normalized-index), 1);
      $non-zero-index: calc($distance-from-guess-card + 0.33);
      $fraction: calc(1 / $non-zero-index);
      $clamped-fraction: clamp(0, $fraction, 1);
      $scale: clamp(0, $clamped-fraction, 1);
      $index-offset: calc($scale - 1);
      $card-offset: calc($index-offset * 50 * $direction);

      $size-scale: calc($scale * 0.5 + 0.5);
      $reverse-size-scale: calc(1 - $size-scale);
      $card-scaled-box-shadow-space: calc(var(--card-box-shadow-space) * $size-scale);

      $card-vertical-offset: calc($card-offset * 1cqh - $direction * $card-height * 0.33);
      $card-scaled-height: calc($card-height * $size-scale);
      $card-scaling-height-loss: calc($reverse-size-scale * $card-height);
      $card-at-top: calc(0% - $card-scaling-height-loss * 0.5 + $card-scaled-box-shadow-space);
      $card-at-vertical-center: calc(50cqh - ($card-scaled-height + $card-scaling-height-loss) * 0.5);
      $card-at-bottom: calc(100cqh - ($card-height - $card-scaling-height-loss * 0.5 + $card-scaled-box-shadow-space));
      $vertical-transform: translateY(clamp($card-at-top, calc($card-at-vertical-center + $card-vertical-offset), $card-at-bottom));
      
      $card-horizontal-offset: calc($card-offset * 1cqw - $direction * $card-width * 0.33);
      $card-scaled-width: calc($card-width * $size-scale);
      $card-scaling-width-loss: calc($reverse-size-scale * $card-width);
      $card-at-left: calc(0% - $card-scaling-width-loss * 0.5 + $card-scaled-box-shadow-space);
      $card-at-horizontal-center: calc(50cqw - ($card-scaled-width + $card-scaling-width-loss) * 0.5);
      $card-at-right: calc(100cqw - ($card-width - $card-scaling-width-loss * 0.5 + $card-scaled-box-shadow-space));
      $horizontal-transform: translateX(clamp($card-at-left, calc($card-at-horizontal-center - $card-horizontal-offset), $card-at-right));

      transform: $vertical-transform $horizontal-transform scale($size-scale);
      
      list-style: none;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      gap: 0.2rem;
      padding: var(--card-padding);
      border-radius: var(--card-padding);
      box-shadow: 0 0 var(--card-box-shadow-space) rgba(0, 0, 0, 0.25);
      background-color: var(--white);
      overflow: hidden;

      transition:
        transform var(--animation-speed-quick),
        z-index var(--animation-speed-quick),
        opacity var(--animation-speed-quick),
        border var(--animation-speed-quick),
        color var(--animation-speed-quick);

      p, h2 {
        max-width: 100%;
        margin: 0;
        overflow: hidden;
        text-overflow: ellipsis;
        text-align: center;
      }

      p {
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        font-size: var(--card-font-size);

        &.track {
          font-size: calc(var(--card-font-size) * 1.5);
        }
      }

      h2 {
        font-size: calc(var(--card-font-size) * 3.75);
      }

      &.guess {
        opacity: 50%;
        border: 2px dashed var(--gray-dark);
      }
    }
  }
</style>
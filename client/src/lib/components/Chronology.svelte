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
      --card-min-height: 8rem;
      --card-scale-height: 50cqh;
      --card-max-height: 16rem;
      --card-font-size: 12cqw;
      --card-aspect-ratio: 1 / 1.25;
      --card-box-shadow-space: 1rem;
      --card-box-shadow-margin: 5px;
      --default-vertical-distance: 15cqh;
      --default-horizontal-distance: 15cqw;
      --distance-step: 0.125;

      $distance-from-guess-card: max(var(--normalized-index), -1 * var(--normalized-index));
      $card-height: clamp(var(--card-min-height), var(--card-scale-height), var(--card-max-height));

      display: flex;
      box-sizing: border-box;
      container-type: size;
      container-name: card;

      aspect-ratio: var(--card-aspect-ratio);
      height: $card-height;
      z-index: calc(var(--card-level) - $distance-from-guess-card);

      position: absolute;
      top: calc(0% + var(--card-box-shadow-space));
      left: calc(0% + var(--card-box-shadow-space));

      $card-at-vertical-center: calc(50cqh - $card-height / 2 - var(--card-box-shadow-space));
      $card-at-bottom: calc(100cqh - $card-height - var(--card-box-shadow-space) * 2);
      $card-width: calc($card-height * var(--card-aspect-ratio));
      $card-at-horizontal-center: calc(50cqw - $card-width / 2 - var(--card-box-shadow-space));
      $card-at-right: calc(100cqw - $card-width - var(--card-box-shadow-space) * 2);
      $direction: clamp(-1, var(--normalized-index), 1);
      $non-zero-index: calc($distance-from-guess-card + 0.5);
      $fraction: calc(1 / $non-zero-index);
      $clamped-fraction: clamp(0, $fraction, 1);
      $index-offset: calc($clamped-fraction - 1);
      $card-offset: calc($index-offset * 50 * $direction);
      
      transform:
        translateY(
          clamp(
            0%,
            calc($card-at-vertical-center + $card-offset * 1cqh),
            $card-at-bottom
          )
        )
        translateX(
          clamp(
            0%,
            calc($card-at-horizontal-center - $card-offset * 1cqw),
            $card-at-right
          )
        )
        scale(
          clamp(
            0,
            $clamped-fraction,
            1
          )
        );
      
      list-style: none;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      gap: 0.2rem;
      padding: var(--card-padding);
      border-radius: var(--card-padding);
      box-shadow: 0 0 calc(var(--card-box-shadow-space) - var(--card-box-shadow-margin)) rgba(255, 1, 213, 1);
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
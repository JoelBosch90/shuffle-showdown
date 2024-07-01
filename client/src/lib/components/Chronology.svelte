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
    if (event instanceof MouseEvent) {
      return {
        x: event.clientX,
        y: event.clientY,
      };
    }
    return {
      x: event.touches[0].clientX,
      y: event.touches[0].clientY,
    };
  };

  const onMoveStart = (event: MouseEvent | TouchEvent) => {
    const { x: startX, y: startY } = getClientLocation(event);
    const startGuessIndex = guessIndex;

    const onMove = (event: MouseEvent | TouchEvent) => {
      const { x: currentX, y: currentY } = getClientLocation(event);
      const heightPercentageTraveled = (currentX - startX) / container.clientWidth;
      const widthPercentageTraveled = (currentY - startY) / container.clientHeight;

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
    position: relative;
    width: 100%;
    padding: min(50%, 8rem);
    flex-grow: 1;
    box-sizing: border-box;
    overflow: hidden;

    cursor: grab;
    user-select: none;
    -webkit-user-select: none;

    &.disabled {
      .card {
        color: var(--gray-dark);
      }
    }

    .card {
      --card-border-radius: 1rem;
      --normalized-index: 0;
      --centered-index: max(var(--normalized-index), -1 * var(--normalized-index));
      --aspect-ratio: 1 / 1.25;
      --default-vertical-distance: 32cqh;
      --default-horizontal-distance: 35cqw;
      --distance-increase: -0.125;

      display: flex;
      box-sizing: border-box;
      container-type: size;
      container-name: card;
      height: 16rem;
      aspect-ratio: 1 / 1.25;
      z-index: calc(var(--card-level) - var(--centered-index));

      position: absolute;
      top: 50%;
      left: 50%;
      transform:
        translateX(
          calc(
            -50% 
            + max(
              min(
                var(--normalized-index) * var(--default-horizontal-distance)
                - var(--distance-increase) * var(--centered-index) * var(--normalized-index) * var(--default-horizontal-distance),
                2 * var(--default-horizontal-distance)
              ),
              -2 * var(--default-horizontal-distance)
            )
          )
        )
        translateY(
          calc(
            -50%
            - max(
              min(
                var(--normalized-index) * var(--default-vertical-distance)
                - var(--distance-increase) * var(--centered-index) * var(--normalized-index) * var(--default-vertical-distance),
                2 * var(--default-vertical-distance)
              ),
              -2 * var(--default-vertical-distance)
            )
          )
        )
        scale(max(calc(1 - var(--centered-index) * 0.2), 0)
      );
      
      list-style: none;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      gap: 0.2rem;
      padding: var(--card-border-radius);
      border-radius: var(--card-border-radius);
      box-shadow: 0 0 1rem rgba(0, 0, 0, 0.1);
      background-color: var(--white);
      overflow: hidden;

      transition:
        transform var(--animation-speed-quick),
        z-index var(--animation-speed-quick),
        opacity var(--animation-speed-quick),
        border var(--animation-speed-quick),
        color var(--animation-speed-quick);

      .track {
        font-size: 2em;
      }

      p {
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        font-size: 1.25em;
      }

      h2 {
        font-size: 5em;
      }

      p, h2 {
        max-width: 100%;
        margin: 0;
        overflow: hidden;
        text-overflow: ellipsis;
        text-align: center;
      }

      &.guess {
        opacity: 0.5;
        border: 2px dashed var(--gray-dark);
      }

      @container chronology (min-height: 20rem) {
        --default-vertical-distance: 6rem;
      }

      @container chronology (min-width: 20rem) {
        --default-horizontal-distance: 7rem;
      }

      @container chronology (max-height: 18rem) {
        height: 10rem;

        h2 { font-size: 2em; }
        p { font-size: 1em;  }
        .track { font-size: 1.25em; }
      }
    }
  }
</style>
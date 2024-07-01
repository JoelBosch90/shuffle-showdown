<script lang="ts">
  import type { WonTrack } from '$lib/types/WonTrack';
  import type { Artist } from '$lib/types/Artist';
  import type { Answer } from '$lib/types/Answer';

  export let disabled: boolean = false;
  export let wonTracks: WonTrack[] = [];
  export let onSelect: (answer: Answer) => void;

  interface Card {
    releaseYear: string;
    name?: string;
    artists?: string;
  }

  const guessCard: Card = {
    releaseYear: "???",
  }

  const trackToCard = ({ track }: WonTrack) : Card => ({
    releaseYear: track.releaseYear?.toString() ?? '???',
    name: track.name,
    artists: joinArtists(track.artists),
  });

  const joinArtists = (artists: Artist[] = []) => {
    return artists.map((artist) => artist.name).join(', ');
  };

  const moveUp = () => {
    updateGuess(guessIndex - 1);
  };

  const moveDown = () => {
    updateGuess(guessIndex + 1);
  };

  const updateGuess = (newIndex: number) => {
    if (disabled) return;

    guessIndex = Math.min(cards.length - 1, Math.max(newIndex, 0));
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

  const sortWonTracks = (a: WonTrack, b: WonTrack) => {
    const aReleaseYear = a.track.releaseYear ?? 0;
    const bReleaseYear = b.track.releaseYear ?? 0;
    return aReleaseYear - bReleaseYear;
  };

  let trackCards: Card[];
  $: trackCards = wonTracks.toSorted(sortWonTracks).map(trackToCard);

  let guessIndex: number;
  $: guessIndex = Math.ceil(wonTracks.length / 2);

  let cards: Card[];
  $: cards = [...trackCards.slice(0, guessIndex), guessCard, ...trackCards.slice(guessIndex)];
  $: cards, selectAnswer();
</script>

<div class="chronology" class:disabled={disabled}>
  {#each cards as card, cardIndex}
    <div
      role="button"
      tabindex="0"
      class="card"
    >
      {#if cardIndex === guessIndex}
        <button class="top-half" on:click={moveUp}></button>
        <button class="bottom-half" on:click={moveDown}></button>
        <i class="fa-solid fa-arrow-up"></i>
      {/if}

      <h2>{card.releaseYear}</h2>

      {#if card.name}
        <p>{card.name}</p>
      {/if}

      {#if card.artists}
        <p>{card.artists}</p>
      {/if}

      {#if cardIndex === guessIndex}
        <i class="fa-solid fa-arrow-down"></i>
      {/if}
    </div>
  {/each}
</div>

<style lang="scss">
  .chronology {
    display: flex;
    flex-direction: column;
    flex-grow: 1;
    box-sizing: border-box;
    align-items: center;
    gap: 1rem;
    padding: 1rem;

    // Make sure the cards can be scrolled through only vertically.
    overflow-y: auto;
    overflow-x: hidden;

    .card {
      --card-border-radius: 1rem;

      position: relative;
      display: flex;
      box-sizing: border-box;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      width: 100%;
      gap: 0.5rem;
      padding: var(--card-border-radius);
      border-radius: var(--card-border-radius);
      box-shadow: 0 0 1rem rgba(0, 0, 0, 0.1);

      h2 {
        font-size: 4rem;
      }

      p, h2 {
        margin: 0;
      }

      .top-half, .bottom-half {
        box-sizing: border-box;
        position: absolute;
        width: 100%;
        height: 50%;
        left: 0;
        cursor: pointer;
      }

      .top-half {
        border-radius: var(--card-border-radius) var(--card-border-radius) 0 0;
        top: 0;
      }

      .bottom-half {
        border-radius: 0 0 var(--card-border-radius) var(--card-border-radius);
        bottom: 0;
      }
    }

    &.disabled {
      .card {
        cursor: not-allowed;
        filter: invert(50%);
      }
    }
  }
</style>
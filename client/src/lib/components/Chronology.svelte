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
    return bReleaseYear + aReleaseYear;
  };

  let trackCards: Card[];
  $: trackCards = wonTracks.toSorted(sortWonTracks).map(trackToCard);

  let guessIndex: number;
  $: guessIndex = Math.ceil(wonTracks.length / 2);

  let cards: Card[];
  $: cards = [...trackCards.slice(0, guessIndex), ...trackCards.slice(0, guessIndex), guessCard, ...trackCards.slice(guessIndex), ...trackCards.slice(guessIndex)];
  $: cards, selectAnswer();
</script>

<ol class="chronology" class:disabled={disabled}>
  {#each cards as card, cardIndex}
    <li class="card {card.releaseYear === '???' ? 'guess' : ''}" style="--normalized-index: {cardIndex - 2}">
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
    position: relative;
    width: 100%;
    flex-grow: 1;
    padding: 1rem;
    box-sizing: border-box;
    overflow: hidden;

    .selection {
      display: flex;
      gap: 1rem;
    }

    .card {
      --card-border-radius: 1rem;
      --normalized-index: 0;
      --centered-index: max(var(--normalized-index), -1 * var(--normalized-index));

      opacity: min(1, max(0, 1.5 - var(--centered-index) * 0.5));
      aspect-ratio: 1 / 1.25;
      height: calc(16rem - (var(--centered-index) * 3rem));

      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      
      container-type: size;
      list-style: none;
      display: flex;
      box-sizing: border-box;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      gap: 0.2rem;
      padding: var(--card-border-radius);
      border-radius: var(--card-border-radius);
      box-shadow: 0 0 1rem rgba(0, 0, 0, 0.1);
      background-color: var(--white);
      overflow: hidden;

      .track {
        font-size: 16cqmin;
      }

      p {
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        font-size: 12cqmin;
      }

      h2 {
        font-size: 45cqmin;
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
    }

    .ghost {
      opacity: 0.5;
      position: absolute;
    }

    &.disabled {
      .card {
        cursor: not-allowed;
        color: var(--gray-dark);
      }
    }
  }
</style>
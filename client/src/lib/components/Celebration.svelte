<script lang="ts">
  import type { Track } from '$lib/types/Track';
  import type { Player } from '$lib/types/Player';
  import type { Artist } from '$lib/types/Artist';
  import type { GameSessionUpdate } from '$lib/types/GameSessionUpdate';
  import { getCurrentRound} from '$lib/helpers/getCurrentRound';
  import { findPlayerInGameSessionUpdate } from '$lib/helpers/findPlayerInGameSessionUpdate';
  import { Timeout } from '$lib/enums/Timeout';

  type CelebrationUpdate = {
    oldUpdate: GameSessionUpdate | null,
    newUpdate: GameSessionUpdate | null,
    oldMe: Player | null,
    newMe: Player | null,
  };

  let track: Track;
  let player: Player;
  let isOtherPlayer: boolean = true;
  let hasWon: boolean = false;
  let isfinalWin: boolean = false;
	let isCelebrating: boolean = false;
	let celebrationPromise: Promise<void> = Promise.resolve();

  const celebrate = (newTrack: Track, newPlayer: Player, newIsOtherPlayer: boolean, newHasWon: boolean, newIsfinalWin: boolean) => {
    track = newTrack;
    player = newPlayer;
    isOtherPlayer = newIsOtherPlayer;
    hasWon = newHasWon;
    isfinalWin = newIsfinalWin;
  }

  const joinArtists = (artists: Artist[] = []) => {
    return artists.map((artist) => artist.name).join(', ');
  };

  const findWinner = (update: GameSessionUpdate, trackPreviewUrl?: string) => {
    if (!trackPreviewUrl) return undefined;

    for (let player of update?.players ?? []) {
      for (let wonTrack of player?.wonTracks ?? []) {
        if (wonTrack.track.previewUrl === trackPreviewUrl) return player;
      }
    }

    return undefined;
  }

  const sleep = async (ms: number) : Promise<void> => {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }

  const celebrateTrack = async (track: Track, player: Player, isOtherPlayer: boolean, hasWon: boolean, isFinalWin: boolean, millisecondsToSleep?: number) => {
    isCelebrating = true;
    celebrate(track, player, isOtherPlayer, hasWon, isFinalWin);

    if (millisecondsToSleep) {
      celebrationPromise = celebrationPromise?.then(() => sleep(millisecondsToSleep));
      await celebrationPromise;
      isCelebrating = false;
    }
  }

  const celebrateRound = async (update: GameSessionUpdate, roundIndex: number, me: Player | null) => {
    const track = update.rounds.find((round) => round.number === roundIndex)?.track;
    if (!track?.name) return;

    const winner = findWinner(update, track.previewUrl);
    const player = winner ?? findPlayerInGameSessionUpdate(update, update.rounds[roundIndex].playerId);

    const isOtherPlayer = player?.id !== me?.id;
    const hasWon = !!winner;

    const wonTracksCount = player?.wonTracks?.length ?? 0;
    const isFinalWin = hasWon && wonTracksCount >= update.songsToWin;

    if (!track || !player) return;
    await celebrateTrack(track, player, isOtherPlayer, hasWon, isFinalWin, Timeout.CELEBRATION);
  }

  const celebrateStart = async (me: Player) => {
    if (!me) return;

    for (const track of me?.wonTracks ?? []) {
      await celebrateTrack(track.track, me, false, true, false, Timeout.CELEBRATION);
    }
  }

  const celebrateEnd = async (update: GameSessionUpdate, newMe: Player | null) => {
    if (!update || !update.hasFinished) return;

    const lastRound = getCurrentRound(update);
    if (!lastRound) return;

    const winner = findWinner(update, lastRound.track.previewUrl);
    if (!winner) return;

    await celebrateTrack(lastRound.track, winner, winner.id !== newMe?.id, true, true);
  }

  const processCelebrations = async (currentGame: GameSessionUpdate | null, update: GameSessionUpdate | null, newMe: Player | null) => {
    if (!update) return;

    if (update.hasFinished) await celebrateEnd(update, newMe);

    const lastSeenRound = getCurrentRound(currentGame)?.number ?? 0;
    const currentRound = getCurrentRound(update)?.number ?? 0;
    if (lastSeenRound >= currentRound) return;

    for (let celebratedRound = lastSeenRound; celebratedRound < currentRound; celebratedRound++) {
      celebrateRound(update, celebratedRound, newMe);
    }
  }

  export const update = async ({ oldUpdate, newUpdate, oldMe, newMe }: CelebrationUpdate) => {
    if (!newUpdate || !newMe) return;

    if (!oldUpdate && 'rounds' in newUpdate && newUpdate.rounds.length > 0) await celebrateStart(newMe);
    if (!oldUpdate && newUpdate.hasFinished) await celebrateEnd(newUpdate, newMe);

		if (newUpdate?.hasFinished) celebrateEnd(newUpdate, newMe);
		else processCelebrations(oldUpdate, newUpdate, newMe);
  }
</script>

<div class="celebration" class:hidden={!isCelebrating}>
  {#if hasWon}
    <h1>Amazing!</h1>
    <p>{isOtherPlayer ? player?.name + " has" : "You have" } won a new track!</p>
  {:else}
    <h1>Oh no!</h1>
    <p>{isOtherPlayer ? player?.name + " has" : "You have" } not won the track!</p>
  {/if}
  {#if track}
    <div class="track">
      <h2>{track.releaseYear}</h2>

      {#if track.name}
        <p>{track.name}</p>
      {/if}

      {#if track.artists}
        <p>{joinArtists(track.artists)}</p>
      {/if}
    </div>
  {/if}
  {#if isfinalWin}
    <h1>Game Winner!</h1>
    <p>Congratulations! {isOtherPlayer ? player?.name + " has" : "You have" } won the game!</p>
  {/if}
</div>

<style lang="scss">
  .celebration {
    position: absolute;
    inset: 0;

    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    gap: 1rem;

		background-color: var(--white);

    text-align: center;

    h1, h2 {
      font-size: 2em;
      margin: 0;
    }

    p {
      margin: 0;
    }

    .track {
      display: flex;
      box-sizing: border-box;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      gap: 0.5rem;

      padding: 1rem;
      border-radius: 1rem;
      box-shadow: 0 0 1rem rgba(0, 0, 0, 0.1);

      margin: 2rem;

      h2 {
        font-size: 4em;
      }
    }

    &.hidden {
      display: none;
    }
  }
</style>
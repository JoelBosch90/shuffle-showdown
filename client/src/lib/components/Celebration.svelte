<script lang="ts">
    import type { Track } from '$lib/types/Track';
    import type { Player } from '$lib/types/Player';
    import type { Artist } from '$lib/types/Artist';

    let track: Track;
    let player: Player;
    let isOtherPlayer: boolean = true;
    let hasWon: boolean = false;
    let isfinalWin: boolean = false;

    const joinArtists = (artists: Artist[] = []) => {
        return artists.map((artist) => artist.name).join(', ');
    };

    export const celebrate = (newTrack: Track, newPlayer: Player, newIsOtherPlayer: boolean, newHasWon: boolean, newIsfinalWin: boolean) => {
        track = newTrack;
        player = newPlayer;
        isOtherPlayer = newIsOtherPlayer;
        hasWon = newHasWon;
        isfinalWin = newIsfinalWin;
    }
</script>

<div class="celebration">
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
        text-align: center;
        
        .track {
            display: flex;
            box-sizing: border-box;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            width: 100%;
            gap: 0.5rem;
            padding: 1rem;
            border-radius: 1rem;
            box-shadow: 0 0 1rem rgba(0, 0, 0, 0.1);

            h2 {
                font-size: 4rem;
                margin: 0;
            }

            p {
                margin: 0;
            }
        }
    }
</style>
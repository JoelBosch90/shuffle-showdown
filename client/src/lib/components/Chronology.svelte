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
        isCurrentGuess: boolean;
    }

    interface CardSlot {
        card?: Card;
    }

    const guessCard: Card = {
        releaseYear: "???",
        isCurrentGuess: true,
    }

    let slots: CardSlot[];
    $: slots = convertToSlots(wonTracks);

    let isDragging: boolean = false;
    let hoveringOption: number | null = null;

    const weaveArrays = <TypeA, TypeB>(arrayA: TypeA[], arrayB: TypeB[]): (TypeA | TypeB)[] => {
        const newArray: (TypeA | TypeB)[] = [];
        const maxLength = Math.max(arrayA.length, arrayB.length);

        for (let index = 0; index < maxLength; index++) {
            if (arrayA[index]) newArray.push(arrayA[index]);
            if (arrayB[index]) newArray.push(arrayB[index]);
        }

        return newArray;
    };

    const joinArtists = (artists: Artist[] = []) => {
        return artists.map((artist) => artist.name).join(', ');
    };

    const sortTracks = (tracks: WonTrack[] = []) => {
        return tracks.toSorted((a, b) => {
            const aNumber = a?.track?.releaseYear ?? 0;
            const bNumber = b?.track?.releaseYear ?? 0;
            return aNumber - bNumber;
        });
    };

    const convertToSlots = (tracks: WonTrack[] = []) => {
        const ownedTracks: CardSlot[] = sortTracks(tracks)
            .map(({ track }) => ({
                card: {
                    releaseYear: track.releaseYear?.toString() ?? "???",
                    name: track.name,
                    artists: joinArtists(track.artists),
                    isCurrentGuess: false,
                }
            }));
        // This looks funky, but the map makes sure that each placeholder has a unique object.
        const guessPlaceholders: CardSlot[] = Array(ownedTracks.length + 1).fill({}).map(() => ({}));
        guessPlaceholders[Math.ceil(guessPlaceholders.length / 2) - 1] = { card: guessCard };

        const slots = weaveArrays(guessPlaceholders, ownedTracks);
        select(slots);
        return slots;
    };

    const select = (slots: CardSlot[] = []) => {
        const slotIndex = slots.findIndex((slot) => slot.card?.isCurrentGuess);

        const cardBefore = slotIndex > 1 ? slots[slotIndex - 1].card : undefined;
        const cardAfter = slotIndex < slots.length - 2 ? slots[slotIndex + 1].card : undefined;

        onSelect({
            afterReleaseYear: cardBefore ? parseInt(cardBefore?.releaseYear ?? "") : undefined,
            beforeReleaseYear: cardAfter ? parseInt(cardAfter?.releaseYear ?? "") : undefined,
        });
    }

    const drag = () => {
        isDragging = true;
    };

    const drop = (event: DragEvent, slotIndex: number) => {
        event.preventDefault();
        isDragging = false;

        const sourceSlot = slots.find((slot) => slot.card?.isCurrentGuess);
        if (sourceSlot) {
            delete sourceSlot.card;   
        }

        slots[slotIndex].card = guessCard;

        select(slots);
    };
</script>

<div class="chronology" class:disabled={disabled}>
    {#each slots as slot, slotIndex}
        {#if slot.card}
            <div
                role="button"
                tabindex="0"
                class="card"
                draggable={slot.card.isCurrentGuess && !disabled}
                on:dragstart={drag}
            >
                <h2>{slot.card.releaseYear}</h2>

                {#if slot.card.name}
                    <p>{slot.card.name}</p>
                {/if}

                {#if slot.card.artists}
                    <p>{slot.card.artists}</p>
                {/if}
            </div>
        {:else if isDragging}
            <div
                role="button"
                tabindex="0"
                class="droppable"
                class:hovering={slotIndex === hoveringOption}
                on:dragenter={() => hoveringOption = slotIndex}
                on:dragleave={() => hoveringOption = null}
                on:drop={(event) => drop(event, slotIndex)}
                on:dragover={(event) => event.preventDefault()}
            />
        {/if}
    {/each}
</div>

<style lang="scss">
    .chronology {
        display: flex;
        flex-direction: column;
        box-sizing: border-box;
        align-items: center;
        gap: 1rem;
        padding: 1rem;
        
        // Make sure the cards can be scrolled through only vertically.
        overflow-y: auto;
        overflow-x: hidden;

        .card {
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

        .droppable {
            display: flex;
            flex-shrink: 0;
            width: 100%;
            height: 6rem;
            background-color: rgba(0, 0, 0, 0.15);
            border-radius: 1rem;
        }

        .hovering {
            background-color: rgba(0, 0, 0, 0.3);
        }

        &.disabled {
            .card {
                cursor: not-allowed;
                filter: invert(50%);
            }
        }
    }
</style>
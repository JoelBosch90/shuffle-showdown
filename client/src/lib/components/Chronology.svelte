<script lang="ts">
    import type { WonTrack } from '$lib/types/WonTrack';
    import type { Artist } from '$lib/types/Artist';
    import type { Answer } from '$lib/types/Answer';

    export let wonTracks: WonTrack[] = [];
    export let onSelect: (answer: Answer) => void;

    interface Card {
        releaseYear: string;
        name?: string;
        artists?: string;
        isGuess: boolean;
    }

    interface CardSlot {
        card?: Card;
    }

    const guessCard: Card = {
        releaseYear: "???",
        isGuess: true,
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

    const convertToSlots = (tracks: WonTrack[]) => {
        const nonPlaceholderTracks: CardSlot[] = tracks.map(({ track }) => ({
            card: {
                releaseYear: track.releaseYear?.toString() ?? "???",
                name: track.name,
                artists: joinArtists(track.artists),
                isGuess: false,
            }
        }));
        const guessPlaceholders: CardSlot[] = Array(nonPlaceholderTracks.length + 1).fill({});
        guessPlaceholders[Math.ceil(guessPlaceholders.length / 2)] = { card: guessCard };

        const slots = weaveArrays(guessPlaceholders, nonPlaceholderTracks);
        select(slots);
        return slots;
    };

    const select = (slots: CardSlot[]) => {
        const slotIndex = slots.findIndex((slot) => slot.card?.isGuess);

        const cardBefore = slotIndex > 1 ? slots[slotIndex - 1].card : undefined;
        const cardAfter = slotIndex < slots.length - 2 ? slots[slotIndex + 1].card : undefined;

        onSelect({
            releaseYearBefore: cardBefore ? parseInt(cardBefore?.releaseYear ?? "") : undefined,
            releaseYearAfter: cardAfter ? parseInt(cardAfter?.releaseYear ?? "") : undefined,
        });
    }

    const drag = () => {
        isDragging = true;
    };

    const drop = (event: DragEvent, slotIndex: number) => {
        event.preventDefault();
        isDragging = false;

        const sourceSlot = slots.find((slot) => slot.card?.isGuess);
        if (sourceSlot) {
            sourceSlot.card = undefined;
        }
        slots[slotIndex].card = guessCard;

        select(slots);
    };
</script>

<div class="chronology">
    {#each slots as slot, slotIndex}
        {#if slot.card}
            <div
                class="chronology-card"
                draggable={slot.card.isGuess}
                on:dragstart={drag}
                role="button"
                tabindex="0"
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
                class="droppable"
                class:hovering={slotIndex === hoveringOption}
                on:dragenter={() => hoveringOption = slotIndex}
                on:dragleave={() => hoveringOption = null}
                on:drop={(event) => drop(event, slotIndex)}
                on:dragover={(event) => event.preventDefault()}
                role="button"
                tabindex="0"
            >
            </div>
        {/if}
    {/each}
</div>

<style lang="scss">
    .chronology {
        display: flex;
        flex-direction: column;
        box-sizing: border-box;
        justify-content: center;
        align-items: center;
        gap: 1rem;
        height: 100%;
    }

    .chronology-card {
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
        width: 100%;
        height: 6rem;
        background-color: rgba(0, 0, 0, 0.15);
        border-radius: 1rem;
    }
    .hovering {
        background-color: rgba(0, 0, 0, 0.3);
    }
</style>
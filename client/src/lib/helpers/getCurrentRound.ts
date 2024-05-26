import type { GameSessionUpdate } from "$lib/types/GameSessionUpdate";
import type { Round } from "$lib/types/Round";

export const getCurrentRound = (update: GameSessionUpdate | null) : Round | null => {
    if (update === null) return null;

    const maxRoundNumber = Math.max(...update.rounds.map((round) => round.number));
    const currentRound = update.rounds.find((round) => round.number === maxRoundNumber);

    return currentRound ?? null;
}
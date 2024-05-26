import type { GameSessionUpdate } from "$lib/types/GameSessionUpdate";

export const findPlayerInGameSessionUpdate = (update: GameSessionUpdate | null, playerId: string | undefined) => update?.players.find((player) => player.id === playerId) ?? null;
    
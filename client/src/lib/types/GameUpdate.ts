import type { Player } from '$lib/types/Player';
import type { Round } from '$lib/types/Round';

export interface GameUpdate {
    id: string;
    configured: boolean;
    hasStarted: boolean;
    hasFinished: boolean;
    songsToWin: number;
    titleRequired: boolean;
    artistRequired: boolean;
    players: Player[];
    rounds: Round[];
}
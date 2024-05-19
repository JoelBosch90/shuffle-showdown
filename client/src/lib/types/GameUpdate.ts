import type { Player } from '$lib/types/Player';
import type { Round } from '$lib/types/Round';

export interface GameUpdate {
    id: string;
    configured: boolean;
    isRunning: boolean;
    songsToWin: number;
    titleRequired: boolean;
    artistRequired: boolean;
    players: Player[];
    rounds: Round[];
}
import type { Player } from '$lib/types/Player';
import type { Playlist } from '$lib/types/Playlist';

export interface GameUpdate {
    id: string;
    configured: boolean;
    isRunning: boolean;
    songsToWin: number;
    titleRequired: boolean;
    artistRequired: boolean;
    players: Player[];
}
import type { Track } from './Track';

export interface Round {
    id: string;
    number: number;
    track: Track;
    playerId: string;
}
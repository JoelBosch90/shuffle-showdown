import type { Track } from "./Track";

export interface WonTrack {
    id: string;
    gameId: string;
    playerId: string;
    track: Track;
}